package graph

// System представляет солнечную систему EVE.
type System struct {
	ID       int
	Name     string
	Security float64
	RegionID int
}

// Graph хранит минимальные данные о карте.
type Graph struct {
	Systems     []System
	Connections [][2]int // каждая пара содержит ID соединённых систем
	Regions     map[int]string
}

// DefaultGraph возвращает небольшой пример графа.
func DefaultGraph() Graph {
	return Graph{
		Systems: []System{
			{ID: 1, Name: "Alpha", Security: 0.5, RegionID: 1},
			{ID: 2, Name: "Beta", Security: 0.6, RegionID: 1},
			{ID: 3, Name: "Gamma", Security: 0.7, RegionID: 1},
		},
		Connections: [][2]int{
			{1, 2},
			{2, 3},
			{1, 3},
		},
		Regions: map[int]string{
			1: "Demo Region",
		},
	}
}
