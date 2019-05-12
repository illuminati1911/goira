package http_test

import (
	"github.com/illuminati1911/goira/internal/models"
	"net/http/cookiejar"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	_authHandler "github.com/illuminati1911/goira/internal/auth/delivery/http"
	_acHandler "github.com/illuminati1911/goira/internal/accontrol/delivery/http"
	"github.com/illuminati1911/goira/testutils"
	_authService "github.com/illuminati1911/goira/internal/auth/service"
	_acservice "github.com/illuminati1911/goira/internal/accontrol/service"
)

var SERVER_URL string

func getFakeState() models.ACState {
	temp := 20
	wind := 0
	mode := 0
	active := false
	return models.ACState{Temperature: &temp, WindLevel: &wind, Mode: &mode, Active: &active}
}

func getServer() *httptest.Server {
	authDB := testutils.NewMockAuthRepository()
	acDB := testutils.NewMockACRepository()
	authServ := _authService.NewAuthService(authDB, "default_pass")
	acServ :=  _acservice.NewACService(acDB, getFakeState(), testutils.NewMockHWInterface())
	mux := http.NewServeMux()
	_authHandler.NewHTTPAuthHandler(authServ, mux)
	_acHandler.NewHTTPACControlHandler(acServ, authServ, mux)
	server := httptest.NewServer(mux)
	SERVER_URL = server.URL
	return server
}

func makeRequest(t *testing.T, c *http.Client, method string, route string, body io.Reader) *http.Response {
	r, err := http.NewRequest(method, SERVER_URL + route, body)
	if err != nil {
		t.Error("Could not create request")
	}
	if body != nil {
		r.Header.Set("Content-type", "application/json")
	}
	resp, err := c.Do(r)
	if err != nil {
		t.Error(err)
	}
	return resp
}

func Login(t *testing.T, c *http.Client) {
	assert := testutils.NewAssert(t)
	rBody, err := json.Marshal(map[string]string{
		"password": "default_pass",
	})
	if err != nil {
		t.Error(err)
	}

	resp := makeRequest(t, c, "POST", "/login", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusOK)
}

func TestGetDefaultState(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()
	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}
	Login(t, c)

	resp := makeRequest(t, c, "GET", "/status", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
	var state models.ACState
	err := json.NewDecoder(resp.Body).Decode(&state)
	assert.Equals(err, nil)
	assert.Equals(*state.Temperature, 20)
	assert.Equals(*state.WindLevel, 0)
	assert.Equals(*state.Mode, 0)
	assert.Equals(*state.Active, false)
}

func TestSetGetState(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()
	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}
	Login(t, c)

	rBody, err := json.Marshal(map[string]interface{}{
		"temp": 28,
	})
	if err != nil {
		t.Error(err)
	}

	resp := makeRequest(t, c, "POST", "/state", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusOK)

	resp = makeRequest(t, c, "GET", "/status", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
	var state models.ACState
	err = json.NewDecoder(resp.Body).Decode(&state)
	assert.Equals(err, nil)
	assert.Equals(*state.Temperature, 28)
	assert.Equals(*state.WindLevel, 0)
	assert.Equals(*state.Mode, 0)
	assert.Equals(*state.Active, false)
}

func TestWrongType(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()
	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}
	Login(t, c)

	rBody, err := json.Marshal(map[string]interface{}{
		"temp": "hello",
	})
	if err != nil {
		t.Error(err)
	}

	resp := makeRequest(t, c, "POST", "/state", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusBadRequest)

	resp = makeRequest(t, c, "GET", "/status", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
	var state models.ACState
	err = json.NewDecoder(resp.Body).Decode(&state)
	assert.Equals(err, nil)
	assert.Equals(*state.Temperature, 20)
	assert.Equals(*state.WindLevel, 0)
	assert.Equals(*state.Mode, 0)
	assert.Equals(*state.Active, false)
}

func TestWrongData(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()
	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}
	Login(t, c)

	rBody, err := json.Marshal(map[string]interface{}{
		"tempdsfd33": "hel23423lo",
		"lol": ":D",
	})
	if err != nil {
		t.Error(err)
	}

	// TODO: Make return 400
	resp := makeRequest(t, c, "POST", "/state", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusOK)

	resp = makeRequest(t, c, "GET", "/status", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
	var state models.ACState
	err = json.NewDecoder(resp.Body).Decode(&state)
	assert.Equals(err, nil)
	assert.Equals(*state.Temperature, 20)
	assert.Equals(*state.WindLevel, 0)
	assert.Equals(*state.Mode, 0)
	assert.Equals(*state.Active, false)
}

func TestFullData(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()
	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}
	Login(t, c)

	rBody, err := json.Marshal(map[string]interface{}{
		"temp": 30,
		"active": true,
		"mode": 2,
		"wind": 1,
	})
	if err != nil {
		t.Error(err)
	}

	// TODO: Make return 400
	resp := makeRequest(t, c, "POST", "/state", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusOK)

	resp = makeRequest(t, c, "GET", "/status", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
	var state models.ACState
	err = json.NewDecoder(resp.Body).Decode(&state)
	assert.Equals(err, nil)
	assert.Equals(*state.Temperature, 30)
	assert.Equals(*state.WindLevel, 1)
	assert.Equals(*state.Mode, 2)
	assert.Equals(*state.Active, true)
}