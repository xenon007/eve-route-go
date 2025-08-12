package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tkhamez/eve-route-go/internal/db"
	routepkg "github.com/tkhamez/eve-route-go/internal/route"
)

func TestNewRouteHandler(t *testing.T) {
	planner, _ := routepkg.NewRoute(db.NewMemory(nil, nil, nil), nil, nil)
	router := mux.NewRouter()
	router.HandleFunc("/api/route/{from}/{to}", NewRouteHandler(planner)).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/api/route/Alpha/Gamma", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var resp struct {
		Routes [][]routepkg.Waypoint `json:"routes"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Routes) == 0 {
		t.Fatalf("expected routes")
	}
	if resp.Routes[0][0].SystemName != "Alpha" {
		t.Fatalf("expected start Alpha, got %s", resp.Routes[0][0].SystemName)
	}
	end := resp.Routes[0][len(resp.Routes[0])-1]
	if end.SystemName != "Gamma" {
		t.Fatalf("expected end Gamma, got %s", end.SystemName)
	}
}

func TestNewRouteHandlerNotFound(t *testing.T) {
	planner, _ := routepkg.NewRoute(db.NewMemory(nil, nil, nil), nil, nil)
	router := mux.NewRouter()
	router.HandleFunc("/api/route/{from}/{to}", NewRouteHandler(planner)).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/api/route/Unknown/Gamma", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
