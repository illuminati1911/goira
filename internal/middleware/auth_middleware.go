package middleware

import (
	"net/http"

	"github.com/illuminati1911/goira/internal/auth"
)

// AuthMiddleware provides HTTP middleware to intercept requests and check
// for the status of the session_token before rpoceeding to the actual business
// logic handling.
//
func AuthMiddleware(as auth.Service) func(f http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(f http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !as.IsAccessTokenValid(c.Value) {
				w.WriteHeader(http.StatusUnauthorized)
			}
			f(w, r)
		}
	}
}
