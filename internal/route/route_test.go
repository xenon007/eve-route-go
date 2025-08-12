package route

import (
	"testing"

	dbstore "github.com/tkhamez/eve-route-go/internal/dbstore"
)

// TestRouteFindAnsiblex проверяет, что Ansiblex рассматривается как альтернативный маршрут.
func TestRouteFindAnsiblex(t *testing.T) {
	ansiblexes := []dbstore.Ansiblex{
		{ID: 1, Name: "Alpha » Gamma - Gate1", SolarSystemID: 1},
		{ID: 2, Name: "Gamma » Alpha - Gate2", SolarSystemID: 3},
	}
	store := dbstore.NewMemory(ansiblexes, nil, nil)
	r, err := NewRoute(store, nil, nil)
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
	temps := []dbstore.TemporaryConnection{
		{System1ID: 1, System2ID: 3},
	}
	store := dbstore.NewMemory(nil, temps, nil)
	r, err := NewRoute(store, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths := r.Find("Alpha", "Gamma")
	if len(paths) != 2 {
		t.Fatalf("ожидались два маршрута, получено %d", len(paths))
	}
	if paths[0][0].ConnectionType == nil || *paths[0][0].ConnectionType != TypeStargate {
		t.Fatalf("первый маршрут должен идти через Stargate")
	}
	if paths[1][0].ConnectionType == nil || *paths[1][0].ConnectionType != TypeTemporary {
		t.Fatalf("второй маршрут должен использовать временное соединение")
	}
}
