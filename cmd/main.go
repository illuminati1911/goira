package main

import (
	"log"
	"os"
	"time"

	accrepo "github.com/illuminati1911/goira/internal/accontrol/repository"

	"github.com/boltdb/bolt"
)

const (
	DBName    string        = "goira.db"
	DBBucket  string        = "accbucket"
	DBMode    os.FileMode   = 0600
	DBTimeout time.Duration = 1 * time.Second
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
	accrepo.NewBoltRepository(db, DBBucket)
}
