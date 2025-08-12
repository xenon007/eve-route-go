package jumps

import (
	"github.com/tkhamez/eve-route-go/internal/graph"
	"testing"
)

func TestCalculator_Between(t *testing.T) {
	t.Run("direct connection", func(t *testing.T) {
		c := NewCalculator(graph.DefaultGraph())
		jumps, err := c.Between("Alpha", "Beta")
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if jumps != 1 {
			t.Fatalf("ожидалось 1 прыжок, получено %d", jumps)
		}
	})

	t.Run("multiple hops", func(t *testing.T) {
		g := graph.Graph{
			Systems: []graph.System{
				{ID: 1, Name: "A", Security: 0, RegionID: 1},
				{ID: 2, Name: "B", Security: 0, RegionID: 1},
				{ID: 3, Name: "C", Security: 0, RegionID: 1},
				{ID: 4, Name: "D", Security: 0, RegionID: 1},
			},
			Connections: [][2]int{{1, 2}, {2, 3}, {3, 4}},
			Regions:     map[int]string{1: "R"},
		}
		c := NewCalculator(g)
		jumps, err := c.Between("A", "D")
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if jumps != 3 {
			t.Fatalf("ожидалось 3 прыжка, получено %d", jumps)
		}
	})

	t.Run("unreachable", func(t *testing.T) {
		g := graph.Graph{
			Systems: []graph.System{
				{ID: 1, Name: "A", Security: 0, RegionID: 1},
				{ID: 2, Name: "B", Security: 0, RegionID: 1},
				{ID: 3, Name: "C", Security: 0, RegionID: 1},
			},
			Connections: [][2]int{{1, 2}},
			Regions:     map[int]string{1: "R"},
		}
		c := NewCalculator(g)
		if _, err := c.Between("A", "C"); err == nil {
			t.Fatalf("ожидалась ошибка для недостижимого маршрута")
		}
	})
}
