package service

import (
	"github.com/illuminati1911/goira/internal/accontrol"
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
func NewACService(repo accontrol.Repository) accontrol.Service {
	return &ACService{repo: repo}
}

// SetTemperature sets the output temperature of the AC. From 16c to 30c.
//
func (acs *ACService) SetTemperature(temp int) error {
	// Call RPi here
	//...

	// if success, save state to DB
	state, err := acs.repo.GetCurrentState()
	if err != nil {
		return err
	}
	state.Temperature = temp
	return acs.repo.SetState(state)
}

// SetWindLevel sets wind level for the AC.
//
func (acs *ACService) SetWindLevel(level int) error {
	// Call RPi here
	//...

	// if success, save state to DB
	state, err := acs.repo.GetCurrentState()
	if err != nil {
		return err
	}
	state.WindLevel = level
	return acs.repo.SetState(state)
}

// TurnOn turns on the AC
//
func (acs *ACService) TurnOn() error {
	// Call RPi here
	//...

	// if success, save state to DB
	state, err := acs.repo.GetCurrentState()
	if err != nil {
		return err
	}
	state.Active = true
	return acs.repo.SetState(state)
}

// TurnOff turns off the AC
//
func (acs *ACService) TurnOff() error {
	// Call RPi here
	//...

	// if success, save state to DB
	state, err := acs.repo.GetCurrentState()
	if err != nil {
		return err
	}
	state.Active = false
	return acs.repo.SetState(state)
}
