package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/illuminati1911/goira/internal/middleware"

	"github.com/illuminati1911/goira/internal/models"

	"github.com/illuminati1911/goira/internal/auth"
)

// HTTPAuthHandler handles authentication related routes with
// any type implementing auth.Service interface
//
type HTTPAuthHandler struct {
	as auth.Service
}

// NewHTTPAuthHandler creates instance of HTTPAuthHandler and sets
// authentication related routes.
//
func NewHTTPAuthHandler(as auth.Service) {
	handler := &HTTPAuthHandler{
		as,
	}
	requireAuth := middleware.AuthMiddleware(as)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/test", requireAuth(handler.Test))
}

// Test is for testing authentication.
// Development use only.
//
func (h *HTTPAuthHandler) Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Valid")
}

// Login is a handler for getting session token
// with given password.
//
func (h *HTTPAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.as.RequestAccessToken(creds.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO: Move the expiration to service.
	//
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(5 * time.Minute),
	})
}
