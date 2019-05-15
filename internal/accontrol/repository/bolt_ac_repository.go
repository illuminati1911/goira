package repository

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// BoltACRepository contains reference to BoltDB instance and bucket name
//
type BoltACRepository struct {
	db     *bolt.DB
	bucket string
}

const (
	StateKey = "state"
)

// NewBoltACRepository returns instance of the BoltRepository implementing
// repository interface.
//
func NewBoltACRepository(db *bolt.DB, bucket string) accontrol.Repository {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return &BoltACRepository{db, bucket}
}

// GetState returns the  state of the system from the repository
//
func (b *BoltACRepository) GetState() (models.ACState, error) {
	var state models.ACState
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(StateKey))
		return json.Unmarshal(v, &state)
	})
	return state, err
}

// SetState stores the state of the system to the repository
//
func (b *BoltACRepository) SetState(state models.ACState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		bytes, err := json.Marshal(state)
		if err != nil {
			return err
		}
		return b.Put([]byte(StateKey), bytes)
	})
}
