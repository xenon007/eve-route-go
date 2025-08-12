package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/tkhamez/eve-route-go/internal/api"
	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/importer"
	routepkg "github.com/tkhamez/eve-route-go/internal/route"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	ctx := context.Background()
	store := db.NewMemory(nil, nil, nil)
	g, err := store.Graph(ctx)
	if err != nil || len(g.Systems) == 0 {
		client := esi.NewClient(nil, "eve-route")
		g, err = importer.BuildGraph(ctx, client)
		if err != nil {
			log.Fatalf("build graph: %v", err)
		}
		if err := store.SaveGraph(ctx, g); err != nil {
			log.Fatalf("save graph: %v", err)
		}
	}
	planner, err := routepkg.NewRoute(store, g, nil, nil)
	if err != nil {
		log.Fatalf("new route: %v", err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/route/{from}/{to}", api.NewRouteHandler(planner)).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.FS(frontendFS)))
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		addr = ":8080"
	}
	log.Printf("server started on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
