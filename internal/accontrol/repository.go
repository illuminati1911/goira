package accontrol

import (
	"github.com/illuminati1911/goira/internal/models"
)

type Repository interface {
	GetState() (models.ACState, error)
	SetState(models.ACState) error
}
