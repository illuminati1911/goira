package repository

import (
	"errors"
	"log"

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

func (b *BoltAuthRepository) SetPassword(pwd models.Password) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(PasswordKey), []byte(pwd))
	})
}

func (b *BoltAuthRepository) IsTokenValid(tkn models.Token) bool {
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		if b.Get([]byte(tkn)) == nil {
			return errors.New("Token is invalid")
		}
		return nil
	})
	return err == nil
}

func (b *BoltAuthRepository) SetToken(tkn models.Token) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(tkn), []byte(tkn))
	})
}

func (b *BoltAuthRepository) DeleteToken(tkn models.Token) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.bucket))
		return b.Put([]byte(TokenKey), []byte(tkn))
	})
}
