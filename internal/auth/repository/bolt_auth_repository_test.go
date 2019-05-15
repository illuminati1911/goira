package repository_test

import (
	"github.com/boltdb/bolt"
	"github.com/illuminati1911/goira/internal/auth"
	"github.com/illuminati1911/goira/internal/auth/repository"
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/testutils"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	DBName    string        = "goira_test.db"
	DBMode    os.FileMode   = 0600
	DBTimeout time.Duration = 1 * time.Second
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

func getRepo() (*bolt.DB, auth.Repository) {
	db := initbolt()
	return db, repository.NewBoltAuthRepository(db, "test_bucket")
}

func TestIfExistingDBIsReused(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	repo.SetPassword("test_password")
	db.Close()
	db2, repo2 := getRepo()
	pwd, _ := repo2.GetPassword()
	assert.Equals(pwd, "test_password")
	db2.Close()
}

func TestPasswordStorage(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	repo.SetPassword("123abc")
	passw, _ := repo.GetPassword()
	assert.Equals(passw, "123abc")
}

func TestTokenStorage(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	currentTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	tkn := models.Token{Name: "token_name", Value: "token_value", Expires: currentTime}
	repo.SetToken(tkn)
	tkn, _ = repo.GetToken("token_value")
	assert.Equals(tkn.Name, "token_name")
	assert.Equals(tkn.Value, "token_value")
	assert.Equals(tkn.Expires, currentTime)
}

func TestDeleteToken(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	testTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	tkn := models.Token{Name: "token_name", Value: "token_value", Expires: testTime}
	repo.SetToken(tkn)
	tkn, _ = repo.GetToken("token_value")
	assert.Equals(tkn.Name, "token_name")
	assert.Equals(tkn.Value, "token_value")
	assert.Equals(tkn.Expires, testTime)
	repo.DeleteToken("token_value")
	_, err := repo.GetToken("token_value")
	assert.NotEquals(err, nil)
}

func TestCleanUp(t *testing.T) {
	defer cleanup()
	assert := testutils.NewAssert(t)
	db, repo := getRepo()
	defer db.Close()
	testTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	tkn := models.Token{Name: "token_name", Value: "token_value", Expires: testTime}
	repo.SetToken(tkn)
	tkn, _ = repo.GetToken("token_value")
	assert.Equals(tkn.Name, "token_name")
	assert.Equals(tkn.Value, "token_value")
	assert.Equals(tkn.Expires, testTime)
	repo.CleanUp()
	_, err := repo.GetToken("token_value")
	assert.NotEquals(err, nil)
}
