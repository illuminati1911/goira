package main

import (
	"log"
	"os"
	"time"

	_accrepo "github.com/illuminati1911/goira/internal/accontrol/repository"
	_acservice "github.com/illuminati1911/goira/internal/accontrol/service"
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
	dbAC := _accrepo.NewBoltACRepository(db, DBACBucket)
	dbAuth := _authrepo.NewBoltAuthRepository(db, DBAuthBucket)
	serviceAC := _acservice.NewACService(dbAC)
	serviceAuth := _authservice.NewAuthService(dbAuth, "dev_pwd")
	println(serviceAC)
	print(serviceAuth)
}
