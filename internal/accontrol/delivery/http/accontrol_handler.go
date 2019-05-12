package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/illuminati1911/goira/internal/models"

	mw "github.com/illuminati1911/goira/internal/middleware"

	"github.com/illuminati1911/goira/internal/accontrol"
	"github.com/illuminati1911/goira/internal/auth"
)

// HTTPACControlHandler handles AC controle related routes
// with any type implementing auth.Service interface
//
type HTTPACControlHandler struct {
	as accontrol.Service
}

// NewHTTPACControlHandler creates instance of HTTPACControlHandlerr and sets
// AC control related routes.
//
func NewHTTPACControlHandler(as accontrol.Service, authS auth.Service, mux *http.ServeMux) *HTTPACControlHandler {
	handler := &HTTPACControlHandler{
		as,
	}
	requireAuth := mw.AuthMiddleware(authS)
	requireAuthPost := mw.Join(requireAuth, mw.PostOnly)
	requireAuthGet := mw.Join(requireAuth, mw.GetOnly)
	mux.HandleFunc("/status", requireAuthGet(handler.GetState))
	mux.HandleFunc("/state", requireAuthPost(handler.SetState))
	return handler
}

func (h *HTTPACControlHandler) GetState(w http.ResponseWriter, r *http.Request) {
	status, err := h.as.GetState()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jStatus, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jStatus)
}

func (h *HTTPACControlHandler) SetState(w http.ResponseWriter, r *http.Request) {
	var state models.ACState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.as.SetState(state) != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
