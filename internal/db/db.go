package db

import "context"

// Ansiblex represents Ansiblex gate data.
type Ansiblex struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	SolarSystemID int    `json:"solarSystemID"`
	RegionID      *int   `json:"regionID,omitempty"`
}

// TemporaryConnection represents temporary connection between two systems.
type TemporaryConnection struct {
	ID        int64 `json:"id"`
	System1ID int   `json:"system1ID"`
	System2ID int   `json:"system2ID"`
}

// System represents a solar system for capital routes.
type System struct {
	ID   int
	Name string
	X    float64
	Y    float64
	Z    float64
}

// Store describes database operations required by the application.
type Store interface {
	Ansiblexes(ctx context.Context) ([]Ansiblex, error)
	TemporaryConnections(ctx context.Context) ([]TemporaryConnection, error)
	Systems(ctx context.Context) (map[int]System, error)

	CreateAnsiblex(ctx context.Context, a Ansiblex) (int64, error)
	UpdateAnsiblex(ctx context.Context, a Ansiblex) error
	DeleteAnsiblex(ctx context.Context, id int64) error

	CreateTemporaryConnection(ctx context.Context, c TemporaryConnection) (int64, error)
	UpdateTemporaryConnection(ctx context.Context, c TemporaryConnection) error
	DeleteTemporaryConnection(ctx context.Context, id int64) error
}
