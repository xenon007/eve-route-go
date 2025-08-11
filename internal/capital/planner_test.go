package capital

import "testing"

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
