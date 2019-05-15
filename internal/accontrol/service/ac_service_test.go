package service_test

import (
	"github.com/illuminati1911/goira/internal/accontrol/service"
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/testutils"
	"testing"
)

func getFakeState() models.ACState {
	temp := 20
	wind := 0
	mode := 0
	active := false
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

func getFakeState2() models.ACState {
	temp := 19
	wind := 1
	mode := 0
	active := true
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

func getFakeState3() models.ACState {
	temp := 30
	return models.ACState{Temperature: &temp, WindLevel: nil, Mode: nil, Active: nil}
}

func TestDefaultModel(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockACRepository()
	state := getFakeState()
	s := service.NewACService(mockDB, state, testutils.NewMockHWInterface())
	gState, err := s.GetState()
	assert.Equals(err, nil)
	assert.Equals(*gState.Temperature, *state.Temperature)
	assert.Equals(*gState.WindLevel, *state.WindLevel)
	assert.Equals(*gState.Mode, *state.Mode)
	assert.Equals(*gState.Active, *state.Active)
}

func TestSetGetModel(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockACRepository()
	state := getFakeState()
	state2 := getFakeState2()
	s := service.NewACService(mockDB, state, testutils.NewMockHWInterface())
	err := s.SetState(state2)
	assert.Equals(err, nil)
	gState, err := s.GetState()
	assert.Equals(err, nil)
	assert.Equals(*gState.Temperature, *state2.Temperature)
	assert.Equals(*gState.WindLevel, *state2.WindLevel)
	assert.Equals(*gState.Mode, *state2.Mode)
	assert.Equals(*gState.Active, *state2.Active)
}

func TestMergingStates(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockACRepository()
	state := getFakeState()
	state3 := getFakeState3()
	s := service.NewACService(mockDB, state, testutils.NewMockHWInterface())
	err := s.SetState(state3)
	assert.Equals(err, nil)
	gState, err := s.GetState()
	assert.Equals(err, nil)
	assert.Equals(*gState.Temperature, *state3.Temperature)
	assert.Equals(*gState.WindLevel, *state.WindLevel)
	assert.Equals(*gState.Mode, *state.Mode)
	assert.Equals(*gState.Active, *state.Active)
}
