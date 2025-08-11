package capital

import (
	"math"
	"testing"
)

func TestPlan(t *testing.T) {
	p := NewPlanner(DefaultSystems(), 5)
	path, err := p.Plan("Maila", "Todifrauan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(path) != 4 || path[0].Name != "Maila" || path[3].Name != "Todifrauan" {
		t.Fatalf("unexpected path: %v", path)
	}
}

func TestPathDistance(t *testing.T) {
	p := NewPlanner(DefaultSystems(), 5)
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
