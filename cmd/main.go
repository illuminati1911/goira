package main

import (
	"fmt"
	"os/signal"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/illuminati1911/goira/internal/models"
	
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/auth"
	_acHandler "github.com/illuminati1911/goira/internal/accontrol/delivery/http"
	"github.com/illuminati1911/goira/internal/accontrol/hwinterface"
	_accrepo "github.com/illuminati1911/goira/internal/accontrol/repository"
	_acservice "github.com/illuminati1911/goira/internal/accontrol/service"
	_authHandler "github.com/illuminati1911/goira/internal/auth/delivery/http"
	_authrepo "github.com/illuminati1911/goira/internal/auth/repository"
	_authservice "github.com/illuminati1911/goira/internal/auth/service"
	"github.com/illuminati1911/goira/internal/accontrol/mappers"

	"github.com/boltdb/bolt"
)

// App wide constants.
const (
	DBName       string        = "goira.db"
	DBACBucket   string        = "accbucket"
	DBAuthBucket string        = "authbucket"
	DBMode       os.FileMode   = 0600
	DBTimeout    time.Duration = 1 * time.Second
)

func initbolt() *bolt.DB {
	db, err := bolt.Open(DBName, DBMode, &bolt.Options{Timeout: DBTimeout})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func acDefaultState() models.ACState {
	temp := 20
	wind := 0
	mode := 0
	active := false
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

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
	serviceAuth := _authservice.NewAuthService(dbAuth, "dev_pwd")
	serviceAC := _acservice.NewACService(dbAC, acDefaultState(), hardwareInfo())
	// HTTP handlers
	//
	mux := http.NewServeMux()
	_authHandler.NewHTTPAuthHandler(serviceAuth, mux)
	_acHandler.NewHTTPACControlHandler(serviceAC, serviceAuth, mux)
	// CleanUp
	//
	t := tokenCleanUp(dbAuth)
	handleShutdown(t, db)
	// Run HTTP
	//
	log.Fatal(http.ListenAndServe(":8080", mux))
}