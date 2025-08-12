package db

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tkhamez/eve-route-go/internal/graph"
)

// TestStore проверяет сохранение и чтение графа.
func TestStore(t *testing.T) {
	g := graph.DefaultGraph()
	dir := t.TempDir()
	path := filepath.Join(dir, "graph.json")
	if err := StoreGraph(context.Background(), path, g); err != nil {
		t.Fatalf("StoreGraph() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	var got graph.Graph
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(g, got) {
		t.Fatalf("stored graph mismatch: got %+v, want %+v", got, g)
	}
}
