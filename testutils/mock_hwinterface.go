package testutils

import (
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// MockHWInterface is a mock HWInterface for testing purposes
//
type MockHWInterface struct{}

// NewMockHWInterface creates a new instance of MockHWInterface
//
func NewMockHWInterface() accontrol.HWInterface {
	return &MockHWInterface{}
}

// SetState does nothing
//
func (m *MockHWInterface) SetState(models.ACState) error {
	return nil
}