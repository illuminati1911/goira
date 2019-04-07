package hwinterface

import (
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

type GPIOInterface struct {
}

func NewGPIOInterface() accontrol.HWInterface {
	return &GPIOInterface{}
}

func (gpio *GPIOInterface) SetState(newState models.ACState) error {
	return nil
}
