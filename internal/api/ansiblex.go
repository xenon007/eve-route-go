package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Ansiblex represents a permanent jump bridge.
type Ansiblex struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

// TempConnection describes a temporary connection between systems.
type TempConnection struct {
	ID      int       `json:"id"`
	From    string    `json:"from"`
	To      string    `json:"to"`
	Expires time.Time `json:"expires"`
}

// Store keeps Ansiblex and temporary connections in memory.
type Store struct {
	mu       sync.Mutex
	nextID   int
	ansiblex map[int]Ansiblex
	temp     map[int]TempConnection
}

// NewStore creates a new in-memory store.
func NewStore() *Store {
	return &Store{ansiblex: make(map[int]Ansiblex), temp: make(map[int]TempConnection), nextID: 1}
}

// RegisterAnsiblexRoutes registers API routes for Ansiblex and temporary connections.
func RegisterAnsiblexRoutes(r *mux.Router, token string) *Store {
	s := NewStore()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/ansiblex", s.listAnsiblex).Methods("GET")
	api.Handle("/ansiblex", s.auth(token, http.HandlerFunc(s.createAnsiblex))).Methods("POST")
	api.Handle("/ansiblex/{id}", s.auth(token, http.HandlerFunc(s.updateAnsiblex))).Methods("PUT")
	api.Handle("/ansiblex/{id}", s.auth(token, http.HandlerFunc(s.deleteAnsiblex))).Methods("DELETE")

	api.HandleFunc("/temp", s.listTemp).Methods("GET")
	api.Handle("/temp", s.auth(token, http.HandlerFunc(s.createTemp))).Methods("POST")
	api.Handle("/temp/{id}", s.auth(token, http.HandlerFunc(s.updateTemp))).Methods("PUT")
	api.Handle("/temp/{id}", s.auth(token, http.HandlerFunc(s.deleteTemp))).Methods("DELETE")

	return s
}

func (s *Store) auth(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Store) listAnsiblex(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []Ansiblex
	for _, a := range s.ansiblex {
		list = append(list, a)
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Store) createAnsiblex(w http.ResponseWriter, r *http.Request) {
	var a Ansiblex
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	a.ID = s.nextID
	s.nextID++
	s.ansiblex[a.ID] = a
	s.mu.Unlock()
	log.Printf("ansiblex created: %d", a.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

func (s *Store) updateAnsiblex(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var a Ansiblex
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	if _, ok := s.ansiblex[id]; !ok {
		s.mu.Unlock()
		http.NotFound(w, r)
		return
	}
	a.ID = id
	s.ansiblex[id] = a
	s.mu.Unlock()
	log.Printf("ansiblex updated: %d", id)
	_ = json.NewEncoder(w).Encode(a)
}

func (s *Store) deleteAnsiblex(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	delete(s.ansiblex, id)
	s.mu.Unlock()
	log.Printf("ansiblex deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Store) listTemp(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []TempConnection
	for _, t := range s.temp {
		list = append(list, t)
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (s *Store) createTemp(w http.ResponseWriter, r *http.Request) {
	var t TempConnection
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	t.ID = s.nextID
	s.nextID++
	s.temp[t.ID] = t
	s.mu.Unlock()
	log.Printf("temp connection created: %d", t.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(t)
}

func (s *Store) updateTemp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var t TempConnection
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	if _, ok := s.temp[id]; !ok {
		s.mu.Unlock()
		http.NotFound(w, r)
		return
	}
	t.ID = id
	s.temp[id] = t
	s.mu.Unlock()
	log.Printf("temp connection updated: %d", id)
	_ = json.NewEncoder(w).Encode(t)
}

func (s *Store) deleteTemp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	delete(s.temp, id)
	s.mu.Unlock()
	log.Printf("temp connection deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}
