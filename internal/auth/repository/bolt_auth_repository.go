package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/illuminati1911/goira/internal/auth"
	"github.com/illuminati1911/goira/internal/models"
)

type BoltAuthRepository struct {
	db     *bolt.DB
	bucket string
}

const (
	PasswordKey = "passwd"
	TokenKey    = "token"
)

// NewBoltAuthRepository returns instance of the BoltAuthRepository implementing
// repository interface defined in auth package.
//
func NewBoltAuthRepository(db *bolt.DB, bucket string) auth.Repository {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	return &BoltAuthRepository{db, bucket}
}

// GetPassword gets password from the repository
//
func (b *BoltAuthRepository) GetPassword() (models.Password, error) {
	var password string
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(PasswordKey))
		if v == nil {
			return errors.New("Password has not been set")
		}
		password = string(v)
		return nil
	})
	return password, err
}

// SetPassword sets password to the repository
//
func (b *BoltAuthRepository) SetPassword(pwd models.Password) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(PasswordKey), []byte(pwd))
	})
}

// GetToken gets token from the repository with token value as a key
//
func (b *BoltAuthRepository) GetToken(tknvalue string) (models.Token, error) {
	var tkn models.Token
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		v := b.Get([]byte(tknvalue))
		if v == nil {
			return errors.New("Token not found")
		}
		return json.Unmarshal(v, &tkn)
	})
	return tkn, err
}

// SetToken sets token to the repository with token value as a key
//
func (b *BoltAuthRepository) SetToken(tkn models.Token) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		bytes, err := json.Marshal(tkn)
		if err != nil {
			return err
		}
		return b.Put([]byte(tkn.Value), bytes)
	})
}

// DeleteToken deletes token from the repository with token value as a key
//
func (b *BoltAuthRepository) DeleteToken(tknvalue string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Delete([]byte(tknvalue))
	})
}

// CleanUp iterates through the repository and deletes expired tokens
//
func (b *BoltAuthRepository) CleanUp() {
	b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		c := b.Cursor()
		var toDelete [][]byte
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tkn models.Token
			if json.Unmarshal(v, &tkn) != nil {
				toDelete = append(toDelete, k)
				continue
			}
			fmt.Printf("key=%s, value=%s\n", k, tkn.Expires)
			if tkn.Expires.Before(time.Now()) {
				toDelete = append(toDelete, k)
			}
		}
		for _, key := range toDelete {
			b.Delete(key)
		}
		return nil
	})
}
