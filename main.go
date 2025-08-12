package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tkhamez/eve-route-go/internal/capital"
	"github.com/tkhamez/eve-route-go/internal/config"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Println("DATABASE_URL is not set")
	}

	r := mux.NewRouter()

	// API endpoint for capital jump planner
	p := capital.NewPlanner(capital.DefaultSystems(), 5)
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

	addr := ":" + cfg.Port
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
