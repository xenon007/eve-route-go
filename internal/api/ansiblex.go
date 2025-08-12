package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tkhamez/eve-route-go/internal/db"
)

// Handler provides HTTP handlers for Ansiblex and temporary connections.
type Handler struct {
	store db.Store
	token string
}

// RegisterAnsiblexRoutes registers API routes for Ansiblex and temporary connections.
func RegisterAnsiblexRoutes(r *mux.Router, token string, store db.Store) *Handler {
	h := &Handler{store: store, token: token}
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/ansiblex", h.listAnsiblex).Methods("GET")
	api.Handle("/ansiblex", h.auth(http.HandlerFunc(h.createAnsiblex))).Methods("POST")
	api.Handle("/ansiblex/{id}", h.auth(http.HandlerFunc(h.updateAnsiblex))).Methods("PUT")
	api.Handle("/ansiblex/{id}", h.auth(http.HandlerFunc(h.deleteAnsiblex))).Methods("DELETE")

	api.HandleFunc("/temp", h.listTemp).Methods("GET")
	api.Handle("/temp", h.auth(http.HandlerFunc(h.createTemp))).Methods("POST")
	api.Handle("/temp/{id}", h.auth(http.HandlerFunc(h.updateTemp))).Methods("PUT")
	api.Handle("/temp/{id}", h.auth(http.HandlerFunc(h.deleteTemp))).Methods("DELETE")

	return h
}

func (h *Handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+h.token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
