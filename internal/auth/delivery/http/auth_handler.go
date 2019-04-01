package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	mw "github.com/illuminati1911/goira/internal/middleware"
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
	requireAuth := mw.AuthMiddleware(as)
	requireAuthGet := mw.Join(requireAuth, mw.GetOnly)
	http.HandleFunc("/login", mw.PostOnly(handler.Login))
	http.HandleFunc("/test", requireAuthGet(handler.Test))
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

	http.SetCookie(w, &http.Cookie{
		Name:    token.Name,
		Value:   token.Value,
		Expires: token.Expires,
	})
}
