package service_test

import (
	"github.com/illuminati1911/goira/internal/auth/service"
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/internal/utils"
	"github.com/illuminati1911/goira/testutils"
	"testing"
	"time"
)

func TestDefaultPassword(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockAuthRepository()
	service.NewAuthService(mockDB, "default_pass")
	password, _ := mockDB.GetPassword()
	assert.Equals(password, utils.SHA256("default_pass"))
}

func TestTokenValidation(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockAuthRepository()
	serv := service.NewAuthService(mockDB, "default_pass")
	validToken := models.Token{Name: "token", Value: "valid", Expires: time.Now().Add(5 * time.Hour)}
	invalidToken := models.Token{
		Name:    "token",
		Value:   "invalid",
		Expires: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)}
	mockDB.SetToken(validToken)
	mockDB.SetToken(invalidToken)
	assert.Equals(serv.IsAccessTokenValid(validToken.Value), true)
	assert.Equals(serv.IsAccessTokenValid(invalidToken.Value), false)
}

func TestFetchAccessToken(t *testing.T) {
	assert := testutils.NewAssert(t)
	mockDB := testutils.NewMockAuthRepository()
	serv := service.NewAuthService(mockDB, "default_pass")
	tkn, err := serv.RequestAccessToken("default_pass")
	assert.Equals(err, nil)
	assert.Equals(serv.IsAccessTokenValid(tkn.Value), true)
	tkn2, err := serv.RequestAccessToken("wrong_pass")
	assert.NotEquals(err, nil)
	assert.Equals(serv.IsAccessTokenValid(tkn2.Value), false)
}
