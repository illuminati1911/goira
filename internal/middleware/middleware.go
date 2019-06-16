package middleware

import (
	"fmt"
	"net/http"

	"github.com/illuminati1911/goira/internal/auth"
)

// MWFunc is a type of middleware function. Takes  HTTP handlerfunction
// as a parameter and returns one.
//
type MWFunc = func(http.HandlerFunc) http.HandlerFunc

// Join is a function to join multiple middleware functions together
// and it returns a single merged function to be used with the final
// endpoint handler.
//
func Join(middlewares ...MWFunc) MWFunc {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return joinHelper(middlewares, final)
	}
}

// AuthMiddleware provides HTTP middleware to intercept requests and check
// for the status of the session_token before rpoceeding to the actual business
// logic handling.
//
func AuthMiddleware(as auth.Service) MWFunc {
	return func(f http.HandlerFunc) http.HandlerFunc {
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

// Cors sets the Cross-Origin policy to allow all origins.
//
// TODO: Make this default for all requests instead of manual input!
//
func Cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			return
		}
		f(w, r)
	}
}

// PostOnly only allows POST requests to proceed.
//
func PostOnly(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodOnly("POST", w, r, f)
	}
}

// GetOnly only allows GET requests to proceed.
//
func GetOnly(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methodOnly("GET", w, r, f)
	}
}

// joinHelper is a recursive function merging middleware
// functions to each other in priority order: left to right.
//
func joinHelper(middlewares []MWFunc, f http.HandlerFunc) http.HandlerFunc {
	if len(middlewares) == 0 {
		return f
	}
	slicer := len(middlewares) - 1
	return joinHelper(middlewares[:slicer], middlewares[slicer:][0](f))
}

// methodOnly is a helper function for method restrictions.
//
func methodOnly(method string, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == method {
		next(w, r)
		return
	}
	http.Error(w, fmt.Sprintf("%s only", method), http.StatusMethodNotAllowed)
}
