package service

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/models"
)

// ACService is a structure containing all the services and action of the AC system.
// This includes saving state to local db as well as communicating with the Raspberry Pi GPIO.
//
type ACService struct {
	repo accontrol.Repository
	hwif accontrol.HWInterface
	mux  sync.Mutex
}

// NewACService creates new ACService with storage system implementing
// repository interface.
//
func NewACService(repo accontrol.Repository, defaultState models.ACState, gpioif accontrol.HWInterface) accontrol.Service {
	_, err := repo.GetState()
	if err == nil {
		return &ACService{repo: repo, hwif: gpioif}
	}

	if repo.SetState(defaultState) != nil {
		log.Fatal("Could not se default state to ACControl")
	}
	return &ACService{repo: repo, hwif: gpioif}
}

// SetState sends the new state to the Hardware interface for IR transmission
// and if it succees saves it to repository.
//
func (acs *ACService) SetState(newState models.ACState) error {
	acs.mux.Lock()
	defer acs.mux.Unlock()
	state, err := acs.repo.GetState()
	if err != nil {
		return err
	}
	mergedState := merge(state, newState)
	if acs.hwif.SetState(mergedState) != nil {
		fmt.Println("GPIO link failure")
		return errors.New("IR Hardware failure")
	}
	return acs.repo.SetState(mergedState)
}

// GetState returns the  state of the system from the repository
//
func (acs *ACService) GetState() (models.ACState, error) {
	acs.mux.Lock()
	defer acs.mux.Unlock()
	return acs.repo.GetState()
}

// Merge merges two AC states.
//
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
	if new.Mode != nil {
		current.Mode = new.Mode
	}
	return current
}
