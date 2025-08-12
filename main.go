package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tkhamez/eve-route-go/internal/capital"
	"github.com/tkhamez/eve-route-go/internal/db"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	r := mux.NewRouter()

	// API endpoint for capital jump planner
	store := db.NewMemory(nil, nil, capital.DefaultSystems())
	p, err := capital.NewPlanner(store, 5)
	if err != nil {
		log.Fatalf("cannot create planner: %v", err)
	}
	r.HandleFunc("/api/capital", func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		if start == "" || end == "" {
			http.Error(w, "missing start or end", http.StatusBadRequest)
			return
		}
		path, err := p.Plan(start, end)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"route": path})
	}).Methods("GET")

	// serve static frontend
	r.PathPrefix("/").Handler(http.FileServer(http.FS(frontendFS)))

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
