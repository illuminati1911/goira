package service

import (
	"errors"

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
	authService := AuthService{repo: repo}
	/*_, err := repo.GetPassword()
	if err == nil {
		return &authService
	}
	if repo.SetPassword(defaultpwd) != nil {
		log.Fatal("Auth: Can't create default password")
	}*/
	return &authService
}

func (as *AuthService) IsAccessTokenValid(tkn models.Token) bool {
	return as.repo.IsTokenValid(tkn)
}

func (as *AuthService) RequestAccessToken(userpwd models.Password) (models.Token, error) {
	systempwd, err := as.repo.GetPassword()
	if err != nil {
		return "", err
	}
	if systempwd != userpwd {
		return "", errors.New("Wrong password")
	}
	token, err := uuid.NewRandom()
	if err != nil {
		return "", errors.New("System failure")
	}
	return token.String(), as.repo.SetToken(token.String())
}
