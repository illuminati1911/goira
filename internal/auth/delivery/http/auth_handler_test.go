package http_test

import (
	"bytes"
	"encoding/json"
	_authHandler "github.com/illuminati1911/goira/internal/auth/delivery/http"
	"github.com/illuminati1911/goira/internal/auth/service"
	"github.com/illuminati1911/goira/testutils"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

var SERVER_URL string

func getServer() *httptest.Server {
	mockDB := testutils.NewMockAuthRepository()
	serv := service.NewAuthService(mockDB, "default_pass")
	mux := http.NewServeMux()
	_authHandler.NewHTTPAuthHandler(serv, mux)
	server := httptest.NewServer(mux)
	SERVER_URL = server.URL
	return server
}

func makeRequest(t *testing.T, c *http.Client, method string, route string, body io.Reader) *http.Response {
	r, err := http.NewRequest(method, SERVER_URL+route, body)
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

func TestTest(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()

	c := &http.Client{}
	resp := makeRequest(t, c, "GET", "/test", nil)
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)
}

func TestLogin(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()

	rBody, err := json.Marshal(map[string]string{
		"password": "default_pass",
	})
	if err != nil {
		t.Error(err)
	}

	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}

	resp := makeRequest(t, c, "GET", "/test", nil)
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)

	resp = makeRequest(t, c, "POST", "/login", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusOK)

	resp = makeRequest(t, c, "GET", "/test", nil)
	assert.Equals(resp.StatusCode, http.StatusOK)
}

func TestFailedLogin(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()

	rBody, err := json.Marshal(map[string]string{
		"password": "wrong_pass",
	})
	if err != nil {
		t.Error(err)
	}

	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}

	resp := makeRequest(t, c, "GET", "/test", nil)
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)

	resp = makeRequest(t, c, "POST", "/login", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)

	resp = makeRequest(t, c, "GET", "/test", nil)
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)
}

func TestBadRequestLogin(t *testing.T) {
	assert := testutils.NewAssert(t)
	server := getServer()
	defer server.Close()

	rBody, err := json.Marshal(map[string]string{
		"not_password": "wrong_pass",
	})
	if err != nil {
		t.Error(err)
	}

	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: cookieJar}

	resp := makeRequest(t, c, "POST", "/login", bytes.NewBuffer(rBody))
	assert.Equals(resp.StatusCode, http.StatusBadRequest)
}
