package route

import (
	"context"
	"testing"

	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/importer"
)

// TestRouteFindAnsiblex проверяет, что Ansiblex рассматривается как альтернативный маршрут.
func TestRouteFindAnsiblex(t *testing.T) {
	ansiblexes := []db.Ansiblex{
		{ID: 1, Name: "Alpha » Gamma - Gate1", SolarSystemID: 1},
		{ID: 2, Name: "Gamma » Alpha - Gate2", SolarSystemID: 3},
	}
	store := db.NewMemory(ansiblexes, nil, nil)
	g, err := importer.BuildGraph(context.Background(), esiWithDirect{})
	if err != nil {
		t.Fatalf("build graph: %v", err)
	}
	_ = store.SaveGraph(context.Background(), g)
	r, err := NewRoute(store, g, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths := r.Find("Alpha", "Gamma")
	t.Logf("paths: %d", len(paths))
	if len(paths) != 2 {
		t.Fatalf("ожидались два маршрута, получено %d", len(paths))
	}
	if paths[0][0].ConnectionType == nil || *paths[0][0].ConnectionType != TypeStargate {
		t.Fatalf("первый маршрут должен идти через Stargate")
	}
	if paths[1][0].ConnectionType == nil || *paths[1][0].ConnectionType != TypeAnsiblex {
		t.Fatalf("второй маршрут должен использовать Ansiblex")
	}
}

// TestRouteSortTemporary проверяет сортировку по количеству временных соединений.
func TestRouteSortTemporary(t *testing.T) {
	temps := []db.TemporaryConnection{{System1ID: 1, System2ID: 3}}
	store := db.NewMemory(nil, temps, nil)
	g, err := importer.BuildGraph(context.Background(), esiNoDirect{})
	if err != nil {
		t.Fatalf("build graph: %v", err)
	}
	_ = store.SaveGraph(context.Background(), g)
	r, err := NewRoute(store, g, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths := r.Find("Alpha", "Gamma")
	if len(paths) != 1 {
		t.Fatalf("ожидался один маршрут, получено %d", len(paths))
	}
	if paths[0][0].ConnectionType == nil || *paths[0][0].ConnectionType != TypeTemporary {
		t.Fatalf("маршрут должен использовать временное соединение")
	}
}

// esiWithDirect реализует ESIClient с прямым соединением Alpha-Gamma.
type esiWithDirect struct{}

func (esiWithDirect) Systems(ctx context.Context) ([]esi.System, error) {
	return []esi.System{
		{ID: 1, Name: "Alpha", Security: 0.5, RegionID: 10},
		{ID: 2, Name: "Beta", Security: 0.6, RegionID: 10},
		{ID: 3, Name: "Gamma", Security: 0.7, RegionID: 10},
	}, nil
}

func (esiWithDirect) Connections(ctx context.Context, systems []esi.System) ([][2]int32, error) {
	return [][2]int32{{1, 2}, {2, 3}, {1, 3}}, nil
}

func (esiWithDirect) RegionName(ctx context.Context, id int32) (string, error) {
	return "Demo", nil
}

// esiNoDirect реализует ESIClient без прямого соединения Alpha-Gamma.
type esiNoDirect struct{}

func (esiNoDirect) Systems(ctx context.Context) ([]esi.System, error) {
	return []esi.System{
		{ID: 1, Name: "Alpha", Security: 0.5, RegionID: 10},
		{ID: 2, Name: "Beta", Security: 0.6, RegionID: 10},
		{ID: 3, Name: "Gamma", Security: 0.7, RegionID: 10},
	}, nil
}

func (esiNoDirect) Connections(ctx context.Context, systems []esi.System) ([][2]int32, error) {
	return [][2]int32{{1, 2}, {2, 3}}, nil
}

func (esiNoDirect) RegionName(ctx context.Context, id int32) (string, error) {
	return "Demo", nil
}
