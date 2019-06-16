package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/illuminati1911/goira/internal/auth"
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/internal/utils"
)

// AuthService is a service structure containing all
// actions related to authentication.
//
type AuthService struct {
	repo auth.Repository
}

// NewAuthService creates new AuthService with storage system implementing
// repository interface defined in auth package.
//
func NewAuthService(repo auth.Repository, defaultpwd models.Password) auth.Service {
	_, err := repo.GetPassword()
	if err == nil {
		fmt.Println("Using existing repository and password.")
		return &AuthService{repo}
	}
	if repo.SetPassword(utils.SHA256(defaultpwd)) != nil {
		log.Fatal("Auth: Can't set default password")
	}
	fmt.Println("Created a new repository with given password.")
	return &AuthService{repo}
}

// IsAccessTokenValid checks if given access token has expired or not
//
func (as *AuthService) IsAccessTokenValid(tknvalue string) bool {
	token, err := as.repo.GetToken(tknvalue)
	if err != nil {
		return false
	}
	if token.Expires.Before(time.Now()) {
		// TODO: Handle error here, don't just return bool
		//
		as.repo.DeleteToken(tknvalue)
		return false
	}
	return true
}

// RequestAccessToken requests for new access token with password as a parameter
//
func (as *AuthService) RequestAccessToken(userpwd models.Password) (models.Token, error) {
	systempwd, err := as.repo.GetPassword()
	if err != nil {
		return models.Token{}, err
	}
	if systempwd != utils.SHA256(userpwd) {
		return models.Token{}, errors.New("Wrong password")
	}
	uuidtoken, err := uuid.NewRandom()
	if err != nil {
		return models.Token{}, errors.New("System failure")
	}
	token := models.Token{Name: "session_token", Value: uuidtoken.String(), Expires: time.Now().Add(1 * time.Hour)}
	return token, as.repo.SetToken(token)
}
