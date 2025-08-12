package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tkhamez/eve-route-go/internal/db"
)

func (h *Handler) listAnsiblex(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.Ansiblexes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) createAnsiblex(w http.ResponseWriter, r *http.Request) {
	var a db.Ansiblex
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.store.CreateAnsiblex(r.Context(), a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.ID = id
	log.Printf("ansiblex created: %d", id)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

func (h *Handler) updateAnsiblex(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var a db.Ansiblex
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.ID = id
	if err := h.store.UpdateAnsiblex(r.Context(), a); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("ansiblex updated: %d", id)
	_ = json.NewEncoder(w).Encode(a)
}

func (h *Handler) deleteAnsiblex(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteAnsiblex(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("ansiblex deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listTemp(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.TemporaryConnections(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) createTemp(w http.ResponseWriter, r *http.Request) {
	var c db.TemporaryConnection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.store.CreateTemporaryConnection(r.Context(), c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.ID = id
	log.Printf("temp connection created: %d", id)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

func (h *Handler) updateTemp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var c db.TemporaryConnection
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	if err := h.store.UpdateTemporaryConnection(r.Context(), c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("temp connection updated: %d", id)
	_ = json.NewEncoder(w).Encode(c)
}

func (h *Handler) deleteTemp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteTemporaryConnection(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("temp connection deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}
