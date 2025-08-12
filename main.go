package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	"github.com/tkhamez/eve-route-go/internal/auth"
	"github.com/tkhamez/eve-route-go/internal/capital"
	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/config"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
  
	cfg := config.FromEnv()
	if cfg.DatabaseURL == "" {
		log.Println("DATABASE_URL is not set")
	}
  
  
	db, err := sql.Open("sqlite", "tokens.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	store, err := auth.NewTokenStore(db)
	if err != nil {
		log.Fatal(err)
	}

	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
	}
  h := auth.NewHandler(conf, store)

	r := mux.NewRouter()
	r.HandleFunc("/login", h.Login).Methods("GET")
	r.HandleFunc("/callback", h.Callback).Methods("GET")

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

	addr := ":" + cfg.Port
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
