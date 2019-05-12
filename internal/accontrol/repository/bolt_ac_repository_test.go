package repository_test

import (
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/testutils"
	"testing"
	"github.com/illuminati1911/goira/internal/accontrol/repository"
	"github.com/illuminati1911/goira/internal/accontrol"
	"path/filepath"
	"log"
	"github.com/boltdb/bolt"
	"time"
	"os"
)

const (
	DBName       string        = "goira_test.db"
	DBMode       os.FileMode   = 0600
	DBTimeout    time.Duration = 1 * time.Second
)

func initbolt() *bolt.DB {
	db, err := bolt.Open(DBName, DBMode, &bolt.Options{Timeout: DBTimeout})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func cleanup() {
	absPath, _ := filepath.Abs(DBName)
	os.Remove(absPath)
}

func getRepo() (*bolt.DB, accontrol.Repository) {
	db := initbolt()
	return db, repository.NewBoltACRepository(db, "test_bucket")
}

func getFakeState() models.ACState {
	temp := 20
	wind := 0
	mode := 0
	active := false
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

func TestIfExistingDBIsReused(t *testing.T) {
	defer cleanup()
	state := getFakeState()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	repo.SetState(state)
	db.Close()
	db2, repo2 := getRepo()
	retainedState, err := repo2.GetCurrentState()
	assert.Equals(err, nil)
	assert.Equals(*retainedState.Temperature, *state.Temperature)
	assert.Equals(*retainedState.WindLevel, *state.WindLevel)
	assert.Equals(*retainedState.Mode, *state.Mode)
	assert.Equals(*retainedState.Active, *state.Active)
	db2.Close()
}


func TestStateSetGet(t *testing.T) {
	defer cleanup()
	state := getFakeState()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	repo.SetState(state)
	gState, err := repo.GetCurrentState()
	assert.Equals(err, nil)
	assert.Equals(*gState.Temperature, *state.Temperature)
	assert.Equals(*gState.WindLevel, *state.WindLevel)
	assert.Equals(*gState.Mode, *state.Mode)
	assert.Equals(*gState.Active, *state.Active)
}

func TestGetNonExisting(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	_, err := repo.GetCurrentState()
	assert.NotEquals(err, nil)
}