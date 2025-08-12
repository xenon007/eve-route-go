package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/importer"
)

func main() {
	ctx := context.Background()
	client := esi.NewClient(nil, "eve-route-importer")
	g, err := importer.BuildGraph(ctx, client)
	if err != nil {
		log.Fatalf("import failed: %v", err)
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(g)
		return
	}
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer dbConn.Close()
	store := db.NewPostgres(dbConn)
	if err := store.SaveGraph(ctx, g); err != nil {
		log.Fatalf("save graph: %v", err)
	}
	log.Printf("graph imported: %d systems, %d connections", len(g.Systems), len(g.Connections))
}
