package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
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
	rows, err := p.db.QueryContext(ctx, "SELECT id, system1_id, system2_id FROM temporary_connections")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []TemporaryConnection
	for rows.Next() {
		var c TemporaryConnection
		if err := rows.Scan(&c.ID, &c.System1ID, &c.System2ID); err != nil {
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

// CreateAnsiblex inserts a new Ansiblex gate.
func (p *Postgres) CreateAnsiblex(ctx context.Context, a Ansiblex) (int64, error) {
	var id int64
	err := p.db.QueryRowContext(ctx, "INSERT INTO ansiblex (name, solar_system_id, region_id) VALUES ($1,$2,$3) RETURNING id", a.Name, a.SolarSystemID, a.RegionID).Scan(&id)
	return id, err
}

// UpdateAnsiblex updates an existing Ansiblex gate.
func (p *Postgres) UpdateAnsiblex(ctx context.Context, a Ansiblex) error {
	_, err := p.db.ExecContext(ctx, "UPDATE ansiblex SET name=$1, solar_system_id=$2, region_id=$3 WHERE id=$4", a.Name, a.SolarSystemID, a.RegionID, a.ID)
	return err
}

// DeleteAnsiblex removes an Ansiblex gate.
func (p *Postgres) DeleteAnsiblex(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM ansiblex WHERE id=$1", id)
	return err
}

// CreateTemporaryConnection inserts a new temporary connection.
func (p *Postgres) CreateTemporaryConnection(ctx context.Context, c TemporaryConnection) (int64, error) {
	var id int64
	err := p.db.QueryRowContext(ctx, "INSERT INTO temporary_connections (system1_id, system2_id) VALUES ($1,$2) RETURNING id", c.System1ID, c.System2ID).Scan(&id)
	return id, err
}

// UpdateTemporaryConnection updates an existing temporary connection.
func (p *Postgres) UpdateTemporaryConnection(ctx context.Context, c TemporaryConnection) error {
	_, err := p.db.ExecContext(ctx, "UPDATE temporary_connections SET system1_id=$1, system2_id=$2 WHERE id=$3", c.System1ID, c.System2ID, c.ID)
	return err
}

// DeleteTemporaryConnection removes a temporary connection.
func (p *Postgres) DeleteTemporaryConnection(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM temporary_connections WHERE id=$1", id)
	return err
}

// EnsurePostgresConnection pings the database to check connection.
func (p *Postgres) EnsurePostgresConnection(ctx context.Context) {
	if err := p.db.PingContext(ctx); err != nil {
		log.Printf("postgres ping error: %v", err)
	}
}
