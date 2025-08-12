package esi

import (
	"context"
	"log"
	"net/http"

	goesi "github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
)

// Client оборачивает библиотеку goesi и предоставляет удобные методы.
type Client struct {
	api *goesi.APIClient
}

// NewClient создаёт новый ESI клиент.
// httpClient может быть nil, в этом случае используется http.DefaultClient.
// userAgent используется для идентификации приложения в запросах.
func NewClient(httpClient *http.Client, userAgent string) *Client {
	log.Printf("esi: new client UA=%s", userAgent)
	api := goesi.NewAPIClient(httpClient, userAgent)
	return &Client{api: api}
}

// Status возвращает статус сервера EVE Online.
func (c *Client) Status(ctx context.Context) (esi.GetStatusOk, error) {
	log.Print("esi: requesting server status")
	status, _, err := c.api.ESI.StatusApi.GetStatus(ctx, nil)
	return status, err
}

// System содержит данные о солнечной системе и списке её стражевых ворот.
type System struct {
	ID        int32
	Name      string
	Security  float64
	RegionID  int32
	Stargates []int32
}

// Systems загружает сведения о всех системах и их стражевых воротах.
func (c *Client) Systems(ctx context.Context) ([]System, error) {
	log.Print("esi: requesting system list")
	ids, _, err := c.api.ESI.UniverseApi.GetUniverseSystems(ctx, nil)
	if err != nil {
		return nil, err
	}
	systems := make([]System, 0, len(ids))
	for _, id := range ids {
		sys, _, err := c.api.ESI.UniverseApi.GetUniverseSystemsSystemId(ctx, id, nil)
		if err != nil {
			return nil, err
		}
		constel, _, err := c.api.ESI.UniverseApi.GetUniverseConstellationsConstellationId(ctx, sys.ConstellationId, nil)
		if err != nil {
			return nil, err
		}
		systems = append(systems, System{
			ID:        id,
			Name:      sys.Name,
			Security:  float64(sys.SecurityStatus),
			RegionID:  constel.RegionId,
			Stargates: sys.Stargates,
		})
	}
	return systems, nil
}

// Connections возвращает уникальные связи между системами по данным о воротах.
func (c *Client) Connections(ctx context.Context, systems []System) ([][2]int32, error) {
	log.Print("esi: requesting connections")
	seen := make(map[[2]int32]struct{})
	for _, s := range systems {
		for _, gateID := range s.Stargates {
			gate, _, err := c.api.ESI.UniverseApi.GetUniverseStargatesStargateId(ctx, gateID, nil)
			if err != nil {
				return nil, err
			}
			pair := [2]int32{s.ID, gate.Destination.SystemId}
			if pair[0] > pair[1] {
				pair[0], pair[1] = pair[1], pair[0]
			}
			seen[pair] = struct{}{}
		}
	}
	res := make([][2]int32, 0, len(seen))
	for p := range seen {
		res = append(res, p)
	}
	return res, nil
}

// RegionName возвращает имя региона по его идентификатору.
func (c *Client) RegionName(ctx context.Context, id int32) (string, error) {
	log.Printf("esi: requesting region %d", id)
	r, _, err := c.api.ESI.UniverseApi.GetUniverseRegionsRegionId(ctx, id, nil)
	if err != nil {
		return "", err
	}
	return r.Name, nil
}
