package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/internal/utils/httputils"

	"github.com/illuminati1911/goira/internal/accontrol"
	_acHandler "github.com/illuminati1911/goira/internal/accontrol/delivery/http"
	"github.com/illuminati1911/goira/internal/accontrol/hwinterface"
	"github.com/illuminati1911/goira/internal/accontrol/mappers"
	_accrepo "github.com/illuminati1911/goira/internal/accontrol/repository"
	_acservice "github.com/illuminati1911/goira/internal/accontrol/service"
	"github.com/illuminati1911/goira/internal/auth"
	_authHandler "github.com/illuminati1911/goira/internal/auth/delivery/http"
	_authrepo "github.com/illuminati1911/goira/internal/auth/repository"
	_authservice "github.com/illuminati1911/goira/internal/auth/service"

	"github.com/boltdb/bolt"
)

// DB Constants
//
const (
	DBName       string        = "goira.db"
	DBACBucket   string        = "accbucket"
	DBAuthBucket string        = "authbucket"
	DBMode       os.FileMode   = 0600
	DBTimeout    time.Duration = 1 * time.Second
)

// Initialize BoltDB database
//
func initbolt() *bolt.DB {
	db, err := bolt.Open(DBName, DBMode, &bolt.Options{Timeout: DBTimeout})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Default state for the AC which will be used if there is no existing state when
// the system starts.
//
func acDefaultState() models.ACState {
	temp := 20
	wind := 0
	mode := 0
	active := false
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

// Get server port from parameters.
//
func getPort() int {
	portPtr := flag.Int("port", 8080, "server port")
	flag.Parse()
	return *portPtr
}

// Gets password from envvar "GOIRA_PASSWORD". If none exists
// uses "default_password".
//
func getEnvPassword() string {
	p := os.Getenv("GOIRA_PASSWORD")
	if p == "" {
		return "default_password"
	}
	return p
}

// Definition for the type of AC to be used as well as the trasmitting GPIO pin.
//
func hardwareInfo() accontrol.HWInterface {
	return hwinterface.NewGPIOInterface(mappers.NewChangHong(), 27)
}

// Make check and clean up fror expired tokens
// every 1 hour.
func tokenCleanUp(repo auth.Repository) *time.Ticker {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			repo.CleanUp()
		}
	}()
	return ticker
}

// Graceful shutdown for SIGINT.
//
func handleShutdown(t *time.Ticker, b *bolt.DB) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("Shutdown...")
			t.Stop()
			b.Close()
			os.Exit(0)
		}
	}()
}

func main() {
	// Init BoltDB.
	//
	db := initbolt()
	defer db.Close()
	// Repositories
	//
	dbAuth := _authrepo.NewBoltAuthRepository(db, DBAuthBucket)
	dbAuth.CleanUp()
	dbAC := _accrepo.NewBoltACRepository(db, DBACBucket)
	// Services
	//
	serviceAuth := _authservice.NewAuthService(dbAuth, getEnvPassword())
	serviceAC := _acservice.NewACService(dbAC, acDefaultState(), hardwareInfo())
	// HTTP handlers
	//
	mux := http.NewServeMux()
	_authHandler.NewHTTPAuthHandler(serviceAuth, mux)
	_acHandler.NewHTTPACControlHandler(serviceAC, serviceAuth, mux)
	// Serve frontend
	//
	mux.Handle("/", httputils.CatchAllHandler(http.FileServer(http.Dir("./frontend"))))
	// CleanUp
	//
	t := tokenCleanUp(dbAuth)
	handleShutdown(t, db)
	// Run HTTP
	//
	fPort := ":" + strconv.Itoa(getPort())
	log.Fatal(http.ListenAndServe(fPort, mux))
}
