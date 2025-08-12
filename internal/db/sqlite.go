package db

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

// SQLite implements Store using SQLite.
type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite store.
func NewSQLite(db *sql.DB) *SQLite {
	return &SQLite{db: db}
}

// Ansiblexes loads Ansiblex gates from SQLite.
func (s *SQLite) Ansiblexes(ctx context.Context) ([]Ansiblex, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, solar_system_id, region_id FROM ansiblex")
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

// TemporaryConnections loads temporary connections from SQLite.
func (s *SQLite) TemporaryConnections(ctx context.Context) ([]TemporaryConnection, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT system1_id, system2_id FROM temporary_connections")
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

// Systems loads capital systems from SQLite.
func (s *SQLite) Systems(ctx context.Context) (map[int]System, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, x, y, z FROM systems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	systems := make(map[int]System)
	for rows.Next() {
		var sys System
		if err := rows.Scan(&sys.ID, &sys.Name, &sys.X, &sys.Y, &sys.Z); err != nil {
			return nil, err
		}
		systems[sys.ID] = sys
	}
	return systems, rows.Err()
}
