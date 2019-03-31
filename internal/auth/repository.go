package auth

import "github.com/illuminati1911/goira/internal/models"

type Repository interface {
	GetPassword() (models.Password, error)
	SetPassword(models.Password) error
	GetToken(string) (models.Token, error)
	SetToken(models.Token) error
	DeleteToken(string) error
}
