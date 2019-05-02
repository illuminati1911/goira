package accontrol

import (
	"github.com/illuminati1911/goira/internal/models"
)

type Mapper interface {
	MapToProtocolBinaryString(ac *models.ACState) string
}
