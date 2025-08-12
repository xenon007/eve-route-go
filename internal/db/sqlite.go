package db

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/tkhamez/eve-route-go/internal/dbstore"
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
func (s *SQLite) Ansiblexes(ctx context.Context) ([]dbstore.Ansiblex, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, solar_system_id, region_id FROM ansiblex")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []dbstore.Ansiblex
	for rows.Next() {
		var a dbstore.Ansiblex
		if err := rows.Scan(&a.ID, &a.Name, &a.SolarSystemID, &a.RegionID); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

// TemporaryConnections loads temporary connections from SQLite.
func (s *SQLite) TemporaryConnections(ctx context.Context) ([]dbstore.TemporaryConnection, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT system1_id, system2_id FROM temporary_connections")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []dbstore.TemporaryConnection
	for rows.Next() {
		var c dbstore.TemporaryConnection
		if err := rows.Scan(&c.System1ID, &c.System2ID); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, rows.Err()
}

// Systems loads capital systems from SQLite.
func (s *SQLite) Systems(ctx context.Context) (map[int]dbstore.System, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, x, y, z FROM systems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	systems := make(map[int]dbstore.System)
	for rows.Next() {
		var sys dbstore.System
		if err := rows.Scan(&sys.ID, &sys.Name, &sys.X, &sys.Y, &sys.Z); err != nil {
			return nil, err
		}
		systems[sys.ID] = sys
	}
	return systems, rows.Err()
}
