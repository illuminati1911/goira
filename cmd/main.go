package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/illuminati1911/goira/internal/models"

	_acHandler "github.com/illuminati1911/goira/internal/accontrol/delivery/http"
	"github.com/illuminati1911/goira/internal/accontrol/hwinterface"
	_accrepo "github.com/illuminati1911/goira/internal/accontrol/repository"
	_acservice "github.com/illuminati1911/goira/internal/accontrol/service"
	_authHandler "github.com/illuminati1911/goira/internal/auth/delivery/http"
	_authrepo "github.com/illuminati1911/goira/internal/auth/repository"
	_authservice "github.com/illuminati1911/goira/internal/auth/service"

	"github.com/boltdb/bolt"
)

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

func main() {
	db := initbolt()
	defer db.Close()
	dbAuth := _authrepo.NewBoltAuthRepository(db, DBAuthBucket)
	dbAC := _accrepo.NewBoltACRepository(db, DBACBucket)
	serviceAuth := _authservice.NewAuthService(dbAuth, "dev_pwd")
	temp := 20
	wind := 0
	mode := 0
	active := false
	gpio := hwinterface.NewGPIOInterface()
	serviceAC := _acservice.NewACService(dbAC, models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}, gpio)
	_authHandler.NewHTTPAuthHandler(serviceAuth)
	_acHandler.NewHTTPACControlHandler(serviceAC, serviceAuth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
