package capital

import (
	"math"
	"testing"

	dbstore "github.com/tkhamez/eve-route-go/internal/dbstore"
)

func TestPlan(t *testing.T) {
	systems := DefaultSystems()
	store := dbstore.NewMemory(nil, nil, systems)
	p, err := NewPlanner(store, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, err := p.Plan("Maila", "Todifrauan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(path) != 4 || path[0].Name != "Maila" || path[3].Name != "Todifrauan" {
		t.Fatalf("unexpected path: %v", path)
	}
}

func TestPathDistance(t *testing.T) {
	systems := DefaultSystems()
	store := dbstore.NewMemory(nil, nil, systems)
	p, err := NewPlanner(store, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path, err := p.Plan("Maila", "Todifrauan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	total := PathDistance(path)
	var expected float64
	for i := 1; i < len(path); i++ {
		expected += distance(path[i-1], path[i])
	}
	if math.Abs(total-expected) > 1e-9 {
		t.Fatalf("unexpected distance: %.6f != %.6f", total, expected)
	}
}
