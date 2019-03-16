package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	fmt.Println("Hello world!")
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
