package route

import "testing"

func TestRouteFindAnsiblex(t *testing.T) {
	ansiblexes := []MongoAnsiblex{
		{ID: 1, Name: "Alpha » Gamma - Gate1", SolarSystemID: 1},
		{ID: 2, Name: "Gamma » Alpha - Gate2", SolarSystemID: 3},
	}
	r := NewRoute(ansiblexes, nil, nil, nil)
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
