package middleware_test

import (
	"github.com/illuminati1911/goira/internal/models"
	"github.com/illuminati1911/goira/internal/auth/service"
	"github.com/illuminati1911/goira/internal/middleware"
	"fmt"
	"github.com/illuminati1911/goira/testutils"
	"io"
	"testing"
	"net/http"
	"net/http/httptest"
)

var SERVER_URL string

func getServerMux() (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	SERVER_URL = server.URL
	return server, mux
}

func makeRequest(t *testing.T, c *http.Client, method string, route string, body io.Reader, sessionToken models.Token) *http.Response {
	r, err := http.NewRequest(method, SERVER_URL + route, body)
	if err != nil {
		t.Error("Could not create request")
	}
	if body != nil {
		r.Header.Set("Content-type", "application/json")
	}
	if sessionToken.Value != "" {
		cookie := &http.Cookie{Name: sessionToken.Name, Value: sessionToken.Value, Expires: sessionToken.Expires}
		r.AddCookie(cookie)
	}
	resp, err := c.Do(r)
	if err != nil {
		t.Error(err)
	}
	return resp
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Valid")
}

func TestPostOnly(t *testing.T) {
	assert := testutils.NewAssert(t)
	server, mux := getServerMux()
	defer server.Close()
	mux.HandleFunc("/test", middleware.PostOnly(test))
	c := &http.Client{}
	resp := makeRequest(t, c, "GET", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusMethodNotAllowed)
	resp = makeRequest(t, c, "POST", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusOK)
}

func TestGetOnly(t *testing.T) {
	assert := testutils.NewAssert(t)
	server, mux := getServerMux()
	defer server.Close()
	mux.HandleFunc("/test", middleware.GetOnly(test))
	c := &http.Client{}
	resp := makeRequest(t, c, "GET", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusOK)
	resp = makeRequest(t, c, "POST", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusMethodNotAllowed)
}

func TestAuthMiddleware(t *testing.T) {
	assert := testutils.NewAssert(t)
	server, mux := getServerMux()
	defer server.Close()
	mockDB := testutils.NewMockAuthRepository()
	serv := service.NewAuthService(mockDB, "default_pass")
	mux.HandleFunc("/test", middleware.AuthMiddleware(serv)(test))
	tkn, _ := serv.RequestAccessToken("default_pass")
	c := &http.Client{}
	resp := makeRequest(t, c, "GET", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)
	resp = makeRequest(t, c, "POST", "/test", nil, tkn)
	assert.Equals(resp.StatusCode, http.StatusOK)
	resp = makeRequest(t, c, "GET", "/test", nil, tkn)
	assert.Equals(resp.StatusCode, http.StatusOK)
}

func TestJoin(t *testing.T) {
	assert := testutils.NewAssert(t)
	server, mux := getServerMux()
	defer server.Close()
	mockDB := testutils.NewMockAuthRepository()
	serv := service.NewAuthService(mockDB, "default_pass")
	requireAuth := middleware.AuthMiddleware(serv)
	requireAuthGet := middleware.Join(requireAuth, middleware.GetOnly)
	mux.HandleFunc("/test", requireAuthGet(test))
	tkn, _ := serv.RequestAccessToken("default_pass")
	c := &http.Client{}
	resp := makeRequest(t, c, "GET", "/test", nil, models.Token{})
	assert.Equals(resp.StatusCode, http.StatusUnauthorized)
	resp = makeRequest(t, c, "POST", "/test", nil, tkn)
	assert.Equals(resp.StatusCode, http.StatusMethodNotAllowed)
	resp = makeRequest(t, c, "GET", "/test", nil, tkn)
	assert.Equals(resp.StatusCode, http.StatusOK)
}