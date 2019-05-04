package testutils

import (
	"github.com/illuminati1911/goira/internal/auth"
	"time"
	"errors"
	"github.com/illuminati1911/goira/internal/models"
)

type MockAuthRepository struct {
	db     map[string]interface{}
}

func NewMockAuthRepository() auth.Repository {
	db := make(map[string]interface{})
	return &MockAuthRepository{db}
}

func (mr *MockAuthRepository) GetPassword() (models.Password, error) {
	password, ok := mr.db["password"].(string)
	if ok {
		return password, nil
	}
	return password, errors.New("Failed to cast value of key: %s to string")
}

func (mr *MockAuthRepository) SetPassword(pwd models.Password) error {
	mr.db["password"] = pwd
	return nil
}

func (mr *MockAuthRepository) GetToken(tknvalue string) (models.Token, error) {
	tkn, ok := mr.db[tknvalue].(models.Token)
	if ok {
		return tkn, nil
	}
	return tkn, errors.New("Failed to cast value of key: %s to token")
}

func (mr *MockAuthRepository) SetToken(tkn models.Token) error {
	mr.db[tkn.Value] = tkn
	return nil
}

func (mr *MockAuthRepository) DeleteToken(tknvalue string) error {
	delete(mr.db, tknvalue)
	return nil
}

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
