package importer

import (
	"context"
	"log"

	"github.com/tkhamez/eve-route-go/internal/esi"
	"github.com/tkhamez/eve-route-go/internal/graph"
)

// ESIClient описывает методы клиента ESI, необходимые для импорта.
type ESIClient interface {
	Systems(ctx context.Context) ([]esi.System, error)
	Connections(ctx context.Context, systems []esi.System) ([][2]int32, error)
	RegionName(ctx context.Context, id int32) (string, error)
}

// BuildGraph загружает данные из ESI и формирует граф карты.
func BuildGraph(ctx context.Context, c ESIClient) (graph.Graph, error) {
	systems, err := c.Systems(ctx)
	if err != nil {
		return graph.Graph{}, err
	}
	conns, err := c.Connections(ctx, systems)
	if err != nil {
		return graph.Graph{}, err
	}
	gSystems := make([]graph.System, 0, len(systems))
	regions := make(map[int]string)
	for _, s := range systems {
		gSystems = append(gSystems, graph.System{
			ID:       int(s.ID),
			Name:     s.Name,
			Security: s.Security,
			RegionID: int(s.RegionID),
		})
		if _, ok := regions[int(s.RegionID)]; !ok {
			name, err := c.RegionName(ctx, s.RegionID)
			if err != nil {
				return graph.Graph{}, err
			}
			regions[int(s.RegionID)] = name
		}
	}
	gConns := make([][2]int, 0, len(conns))
	for _, p := range conns {
		gConns = append(gConns, [2]int{int(p[0]), int(p[1])})
	}
	log.Printf("importer: imported %d systems and %d connections", len(gSystems), len(gConns))
	return graph.Graph{Systems: gSystems, Connections: gConns, Regions: regions}, nil
}
