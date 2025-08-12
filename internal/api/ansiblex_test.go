package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestAnsiblexAuth(t *testing.T) {
	r := mux.NewRouter()
	RegisterAnsiblexRoutes(r, "token")

	body := bytes.NewBufferString(`{"name":"A","from":"X","to":"Y"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/ansiblex", body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	body = bytes.NewBufferString(`{"name":"A","from":"X","to":"Y"}`)
	req = httptest.NewRequest(http.MethodPost, "/api/ansiblex", body)
	req.Header.Set("Authorization", "Bearer token")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/ansiblex", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var list []Ansiblex
	if err := json.NewDecoder(w.Body).Decode(&list); err != nil || len(list) != 1 {
		t.Fatalf("expected one element, got %v %v", len(list), err)
	}
}

func TestTempAuth(t *testing.T) {
	r := mux.NewRouter()
	RegisterAnsiblexRoutes(r, "token")

	body := bytes.NewBufferString(`{"from":"X","to":"Y"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/temp", body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
