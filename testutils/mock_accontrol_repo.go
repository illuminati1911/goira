package testutils

import (
	"errors"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// MockACRepository is a in-memory mock repository for testing
// purposes
//
type MockACRepository struct {
	state *models.ACState
}

// NewMockACRepository creates a new instance of MockACRepository
//
func NewMockACRepository() accontrol.Repository {
	return &MockACRepository{nil}
}

// GetState returns  state from memory
//
func (m *MockACRepository) GetState() (models.ACState, error) {
	if m.state != nil {
		return *m.state, nil
	}
	return models.ACState{}, errors.New("Failure")
}

// SetState sets state to memory
//
func (m *MockACRepository) SetState(state models.ACState) error {
	m.state = &state
	return nil
}
