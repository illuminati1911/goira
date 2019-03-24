package auth

import "github.com/illuminati1911/goira/internal/models"

type Service interface {
	IsAccessTokenValid(token models.Token) bool
	RequestAccessToken(models.Password) (models.Token, error)
}
