package service

import (
	"log"

	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// ACService is a structure containing all the services and action of the AC system.
// This includes saving state to local db as well as communicating with the Raspberry Pi GPIO.
//
type ACService struct {
	repo accontrol.Repository
}

// NewACService creates new ACService with storage system implementing
// repository interface.
//
func NewACService(repo accontrol.Repository, defaultState models.ACState) accontrol.Service {
	_, err := repo.GetCurrentState()
	if err == nil {
		return &ACService{repo: repo}
	}

	if repo.SetState(defaultState) != nil {
		log.Fatal("Could not se default state to ACControl")
	}
	return &ACService{repo: repo}
}

func (acs *ACService) SetState(newState models.ACState) error {
	state, err := acs.repo.GetCurrentState()
	if err != nil {
		return err
	}
	mergedState := merge(state, newState)
	// Call RPi here
	//...

	// if success, save state to DB
	return acs.repo.SetState(mergedState)
}

func (acs *ACService) GetState() (models.ACState, error) {
	return acs.repo.GetCurrentState()
}

func merge(current models.ACState, new models.ACState) models.ACState {
	if new.Temperature != nil {
		current.Temperature = new.Temperature
	}
	if new.Active != nil {
		current.Active = new.Active
	}
	if new.WindLevel != nil {
		current.WindLevel = new.WindLevel
	}
	return current
}
