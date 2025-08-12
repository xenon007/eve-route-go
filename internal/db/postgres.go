package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"

	"github.com/tkhamez/eve-route-go/internal/graph"
)

// Postgres implements Store using PostgreSQL.
type Postgres struct {
	db *sql.DB
}

// NewPostgres creates a new Postgres store.
func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// Ansiblexes loads Ansiblex gates from PostgreSQL.
func (p *Postgres) Ansiblexes(ctx context.Context) ([]Ansiblex, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, name, solar_system_id, region_id FROM ansiblex")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Ansiblex
	for rows.Next() {
		var a Ansiblex
		if err := rows.Scan(&a.ID, &a.Name, &a.SolarSystemID, &a.RegionID); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

// TemporaryConnections loads temporary connections from PostgreSQL.
func (p *Postgres) TemporaryConnections(ctx context.Context) ([]TemporaryConnection, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT system1_id, system2_id FROM temporary_connections")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []TemporaryConnection
	for rows.Next() {
		var c TemporaryConnection
		if err := rows.Scan(&c.System1ID, &c.System2ID); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, rows.Err()
}

// Systems loads capital systems from PostgreSQL.
func (p *Postgres) Systems(ctx context.Context) (map[int]System, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, name, x, y, z FROM systems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	systems := make(map[int]System)
	for rows.Next() {
		var s System
		if err := rows.Scan(&s.ID, &s.Name, &s.X, &s.Y, &s.Z); err != nil {
			return nil, err
		}
		systems[s.ID] = s
	}
	return systems, rows.Err()
}

// Graph loads the universe graph from PostgreSQL.
func (p *Postgres) Graph(ctx context.Context) (graph.Graph, error) {
	var data []byte
	err := p.db.QueryRowContext(ctx, "SELECT data FROM graph LIMIT 1").Scan(&data)
	if err != nil {
		return graph.Graph{}, err
	}
	var g graph.Graph
	if err := json.Unmarshal(data, &g); err != nil {
		return graph.Graph{}, err
	}
	return g, nil
}

// SaveGraph stores the universe graph in PostgreSQL.
func (p *Postgres) SaveGraph(ctx context.Context, g graph.Graph) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}
	_, err = p.db.ExecContext(ctx, "INSERT INTO graph (id, data) VALUES (1, $1) ON CONFLICT (id) DO UPDATE SET data = EXCLUDED.data", data)
	return err
}

// EnsurePostgresConnection pings the database to check connection.
func (p *Postgres) EnsurePostgresConnection(ctx context.Context) {
	if err := p.db.PingContext(ctx); err != nil {
		log.Printf("postgres ping error: %v", err)
	}
}
