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
