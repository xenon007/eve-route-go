package db

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/tkhamez/eve-route-go/internal/graph"
)

// StoreGraph сохраняет граф в указанном файле в формате JSON.
func StoreGraph(_ context.Context, path string, g graph.Graph) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("db: close file: %v", cErr)
		}
	}()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(g); err != nil {
		return err
	}
	log.Printf("db: stored graph with %d systems to %s", len(g.Systems), path)
	return nil
}
