package testutils

import (
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

type MockHWInterface struct{}

func NewMockHWInterface() accontrol.HWInterface {
	return &MockHWInterface{}
}

func (m *MockHWInterface) SetState(models.ACState) error {
	return nil
}