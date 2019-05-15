package testutils

import (
	"errors"
	"github.com/illuminati1911/goira/internal/auth"
	"github.com/illuminati1911/goira/internal/models"
	"time"
)

// MockAuthRepository is a in-memory mock repository for testing
// purposes
//
type MockAuthRepository struct {
	db map[string]interface{}
}

// NewMockAuthRepository creates a new instance of MockAuthRepository
//
func NewMockAuthRepository() auth.Repository {
	db := make(map[string]interface{})
	return &MockAuthRepository{db}
}

// GetPassword returns password from the in-memory map
//
func (mr *MockAuthRepository) GetPassword() (models.Password, error) {
	password, ok := mr.db["password"].(string)
	if ok {
		return password, nil
	}
	return password, errors.New("Failed to cast value of key: %s to string")
}

// SetPassword sets password to the in-memory map
//
func (mr *MockAuthRepository) SetPassword(pwd models.Password) error {
	mr.db["password"] = pwd
	return nil
}

// GetToken returns full token from corrensponding token value
//
func (mr *MockAuthRepository) GetToken(tknvalue string) (models.Token, error) {
	tkn, ok := mr.db[tknvalue].(models.Token)
	if ok {
		return tkn, nil
	}
	return tkn, errors.New("Failed to cast value of key: %s to token")
}

// SetToken sets token to the map with token value as key
//
func (mr *MockAuthRepository) SetToken(tkn models.Token) error {
	mr.db[tkn.Value] = tkn
	return nil
}

// DeleteToken deletes the token from the map with token value as key
//
func (mr *MockAuthRepository) DeleteToken(tknvalue string) error {
	delete(mr.db, tknvalue)
	return nil
}

// CleanUp iterates through the repository and deletes expired tokens
//
func (mr *MockAuthRepository) CleanUp() {
	var toDelete []models.Token
	for key, value := range mr.db {
		if key == "password" {
			continue
		}
		tkn, ok := value.(models.Token)
		if ok && tkn.Expires.Before(time.Now()) {
			toDelete = append(toDelete, tkn)
		}
	}
	for _, t := range toDelete {
		delete(mr.db, t.Value)
	}
}
