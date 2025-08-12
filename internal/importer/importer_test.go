package importer

import (
	"context"
	"testing"

	"github.com/tkhamez/eve-route-go/internal/esi"
)

type fakeESI struct{}

func (fakeESI) Systems(ctx context.Context) ([]esi.System, error) {
	return []esi.System{
		{ID: 1, Name: "Alpha", Security: 0.5, RegionID: 10},
		{ID: 2, Name: "Beta", Security: 0.6, RegionID: 10},
	}, nil
}

func (fakeESI) Connections(ctx context.Context, systems []esi.System) ([][2]int32, error) {
	return [][2]int32{{1, 2}}, nil
}

func (fakeESI) RegionName(ctx context.Context, id int32) (string, error) {
	return map[int32]string{10: "Demo"}[id], nil
}

func TestBuildGraph(t *testing.T) {
	g, err := BuildGraph(context.Background(), fakeESI{})
	if err != nil {
		t.Fatalf("build graph: %v", err)
	}
	if len(g.Systems) != 2 {
		t.Fatalf("expected 2 systems, got %d", len(g.Systems))
	}
	if len(g.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(g.Connections))
	}
	if g.Regions[10] != "Demo" {
		t.Fatalf("unexpected region map: %#v", g.Regions)
	}
}
