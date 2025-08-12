package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tkhamez/eve-route-go/internal/capital"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	r := mux.NewRouter()

	// API endpoint for capital jump planner
	p := capital.NewPlanner(capital.DefaultSystems(), 5)
	r.HandleFunc("/api/capital", func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		if start == "" || end == "" {
			log.Printf("capital api: missing start or end (start=%q end=%q)", start, end)
			http.Error(w, "missing start or end", http.StatusBadRequest)
			return
		}
		path, err := p.Plan(start, end)
		if err != nil {
			log.Printf("capital api: %v", err)
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
