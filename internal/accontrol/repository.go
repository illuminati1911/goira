package accontrol

import (
	"github.com/illuminati1911/goira/internal/models"
)

type Repository interface {
	GetCurrentState() (models.ACState, error)
	GetPassword() (string, error)
	GetToken() (string, error)
	SetState(models.ACState) error
	SetPassword(string) error
	SetToken(string) error
}
