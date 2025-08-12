package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tkhamez/eve-route-go/internal/db"
)

func TestAnsiblexCRUD(t *testing.T) {
	r := mux.NewRouter()
	store := db.NewMemory(nil, nil, nil)
	RegisterAnsiblexRoutes(r, "token", store)

	body := bytes.NewBufferString(`{"name":"A","solarSystemID":1}`)
	req := httptest.NewRequest(http.MethodPost, "/api/ansiblex", body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	body = bytes.NewBufferString(`{"name":"A","solarSystemID":1}`)
	req = httptest.NewRequest(http.MethodPost, "/api/ansiblex", body)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created db.Ansiblex
	_ = json.NewDecoder(w.Body).Decode(&created)

	body = bytes.NewBufferString(`{"name":"B","solarSystemID":1}`)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/ansiblex/%d", created.ID), body)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/ansiblex", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []db.Ansiblex
	_ = json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 1 || list[0].Name != "B" {
		t.Fatalf("unexpected list %v", list)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/ansiblex/%d", created.ID), nil)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/ansiblex", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	_ = json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %v", list)
	}
}

func TestTempCRUD(t *testing.T) {
	r := mux.NewRouter()
	store := db.NewMemory(nil, nil, nil)
	RegisterAnsiblexRoutes(r, "token", store)

	body := bytes.NewBufferString(`{"system1ID":1,"system2ID":2}`)
	req := httptest.NewRequest(http.MethodPost, "/api/temp", body)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var created db.TemporaryConnection
	_ = json.NewDecoder(w.Body).Decode(&created)

	body = bytes.NewBufferString(`{"system1ID":3,"system2ID":4}`)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/temp/%d", created.ID), body)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/temp", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []db.TemporaryConnection
	_ = json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 1 || list[0].System1ID != 3 {
		t.Fatalf("unexpected list %v", list)
	}

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/temp/%d", created.ID), nil)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/temp", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	_ = json.NewDecoder(w.Body).Decode(&list)
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %v", list)
	}
}
