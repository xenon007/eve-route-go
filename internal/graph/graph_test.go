package graph

import (
	"log"
	"reflect"
	"testing"
)

// TestDefaultGraph проверяет создание стандартного графа.
func TestDefaultGraph(t *testing.T) {
	g := DefaultGraph()
	log.Printf("получен граф: %+v", g)

	if len(g.Systems) != 3 {
		t.Fatalf("ожидалось 3 системы, получено %d", len(g.Systems))
	}
	if g.Systems[0].Name != "Alpha" {
		t.Errorf("ожидалась первая система Alpha, получено %s", g.Systems[0].Name)
	}
	expectedRegions := map[int]string{1: "Demo Region"}
	if !reflect.DeepEqual(g.Regions, expectedRegions) {
		t.Errorf("неверные регионы: %+v", g.Regions)
	}
	if len(g.Connections) != 3 {
		t.Errorf("ожидалось 3 соединения, получено %d", len(g.Connections))
	}
}
