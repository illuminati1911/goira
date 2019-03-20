package repository

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

type BoltRepository struct {
	db     *bolt.DB
	bucket string
}

const (
	StateKey    = "state"
	PasswordKey = "passwd"
	TokenKey    = "token"
)

// NewBoltRepository returns instance of the BoltRepository implementing
// repository interface.
//
func NewBoltRepository(db *bolt.DB, bucket string) accontrol.Repository {
	return &BoltRepository{db, bucket}
}

func (b *BoltRepository) GetCurrentState() (models.ACState, error) {
	var state models.ACState
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(StateKey))
		return json.Unmarshal(v, &state)
	})
	return state, err
}

func (b *BoltRepository) GetPassword() (string, error) {
	var password string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(PasswordKey))
		password = string(v)
		return nil
	})
	return password, err
}

func (b *BoltRepository) GetToken() (string, error) {
	var token string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(PasswordKey))
		token = string(v)
		return nil
	})
	return token, err
}

func (b *BoltRepository) SetState(state models.ACState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		bytes, err := json.Marshal(state)
		if err != nil {
			return err
		}
		return b.Put([]byte(StateKey), bytes)
	})
}

func (b *BoltRepository) SetPassword(pwd string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(PasswordKey), []byte(pwd))
	})
}

func (b *BoltRepository) SetToken(tkn string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(TokenKey), []byte(tkn))
	})
}
