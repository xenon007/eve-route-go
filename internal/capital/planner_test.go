package capital

import "testing"

func TestPlan(t *testing.T) {
	p := NewPlanner(DefaultSystems(), DefaultGraph())
	path, err := p.Plan("Maila", "Todifrauan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(path) != 2 || path[0].Name != "Maila" || path[1].Name != "Todifrauan" {
		t.Fatalf("unexpected path: %v", path)
	}
}
