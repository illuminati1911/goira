package accontrol

import (
	"github.com/illuminati1911/goira/internal/models"
)

type HWInterface interface {
	SetState(models.ACState) error
}
