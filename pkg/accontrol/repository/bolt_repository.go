package accrepo

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/illuminati1911/goira/pkg/accontrol"
	"github.com/illuminati1911/goira/pkg/models"
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
		if err := json.Unmarshal(v, state); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.ACState{}, err
	}
	return state, nil
}

func (b *BoltRepository) GetPassword() (string, error) {
	var password string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(PasswordKey))
		password = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return password, nil
}

func (b *BoltRepository) GetToken() (string, error) {
	var token string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(PasswordKey))
		token = string(v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (b *BoltRepository) SetState(state models.ACState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		bytes, err := json.Marshal(state)
		if err != nil {
			return err
		}
		err = b.Put([]byte(StateKey), bytes)
		return err
	})
}

func (b *BoltRepository) SetPassword(pwd string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		err := b.Put([]byte(PasswordKey), []byte(pwd))
		return err
	})
}

func (b *BoltRepository) SetToken(tkn string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		err := b.Put([]byte(TokenKey), []byte(tkn))
		return err
	})
}
