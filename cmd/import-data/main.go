package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/importer"
)

func main() {
	ctx := context.Background()
	client := esi.NewClient(nil, "eve-route-importer")
	graph, err := importer.BuildGraph(ctx, client)
	if err != nil {
		log.Fatalf("import failed: %v", err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(graph); err != nil {
		log.Fatalf("encode graph: %v", err)
	}
}
