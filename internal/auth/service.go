package auth

import "github.com/illuminati1911/goira/internal/models"

type Service interface {
	IsAccessTokenValid(string) bool
	RequestAccessToken(models.Password) (models.Token, error)
}
