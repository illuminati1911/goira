package accontrol

import (
	"github.com/illuminati1911/goira/internal/models"
)

type Service interface {
	SetState(models.ACState) (models.ACState, error)
	GetState() (models.ACState, error)
}
