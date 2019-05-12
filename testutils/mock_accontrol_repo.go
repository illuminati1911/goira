package testutils

import (
	"errors"
	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

type MockACRepository struct {
	state     *models.ACState
}

func NewMockACRepository() accontrol.Repository {
	return &MockACRepository{nil}
}

func (m *MockACRepository) GetCurrentState() (models.ACState, error) {
	if m.state != nil {
		return *m.state, nil
	}
	return models.ACState{}, errors.New("Failure")
}

func (m *MockACRepository) SetState(state models.ACState) error {
	m.state = &state
	return nil
}