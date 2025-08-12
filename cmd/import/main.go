package main

import (
	"context"
	"log"
	"os"

	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/importer"
)

// main запускает процесс импорта данных ESI и сохраняет граф.
func main() {
	ctx := context.Background()
	client := esi.NewClient(nil, "eve-route-importer")
	log.Println("import: building graph from ESI")
	g, err := importer.BuildGraph(ctx, client)
	if err != nil {
		log.Fatalf("import: build graph: %v", err)
	}
	path := "graph.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if err := db.StoreGraph(ctx, path, g); err != nil {
		log.Fatalf("import: store graph: %v", err)
	}
	log.Printf("import: graph saved to %s", path)
}
