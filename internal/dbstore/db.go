package dbstore

import "context"

// Ansiblex represents Ansiblex gate data.
type Ansiblex struct {
	ID            int64
	Name          string
	SolarSystemID int
	RegionID      *int
}

// TemporaryConnection represents temporary connection between two systems.
type TemporaryConnection struct {
	System1ID int
	System2ID int
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
}
