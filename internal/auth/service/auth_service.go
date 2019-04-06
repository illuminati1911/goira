package service

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/illuminati1911/goira/internal/auth"
	"github.com/illuminati1911/goira/internal/models"
)

type AuthService struct {
	repo auth.Repository
}

// NewAuthService creates new AuthService with storage system implementing
// repository interface defined in auth package.
//
func NewAuthService(repo auth.Repository, defaultpwd models.Password) auth.Service {
	_, err := repo.GetPassword()
	if err == nil {
		return &AuthService{repo: repo}
	}
	if repo.SetPassword(defaultpwd) != nil {
		log.Fatal("Auth: Can't create default password")
	}
	return &AuthService{repo: repo}
}

func (as *AuthService) IsAccessTokenValid(tknvalue string) bool {
	token, err := as.repo.GetToken(tknvalue)
	if err != nil {
		return false
	}
	if token.Expires.Before(time.Now()) {
		// Handle error here, don't just return bool
		//
		as.repo.DeleteToken(tknvalue)
		return false
	}
	return true
}

func (as *AuthService) RequestAccessToken(userpwd models.Password) (models.Token, error) {
	systempwd, err := as.repo.GetPassword()
	if err != nil {
		return models.Token{}, err
	}
	if systempwd != userpwd {
		return models.Token{}, errors.New("Wrong password")
	}
	uuidtoken, err := uuid.NewRandom()
	if err != nil {
		return models.Token{}, errors.New("System failure")
	}
	token := models.Token{Name: "session_token", Value: uuidtoken.String(), Expires: time.Now().Add(5 * time.Minute)}
	return token, as.repo.SetToken(token)
}
