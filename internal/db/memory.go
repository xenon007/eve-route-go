package db

import (
	"context"
	"errors"
)

// Memory provides an in-memory implementation of Store for tests.
type Memory struct {
	ansiblexes      []Ansiblex
	tempConnections []TemporaryConnection
	systems         map[int]System
	nextAnsiblexID  int64
	nextTempID      int64
}

// NewMemory creates a new in-memory store instance.
func NewMemory(ans []Ansiblex, temps []TemporaryConnection, systems map[int]System) *Memory {
	if ans == nil {
		ans = []Ansiblex{}
	}
	if temps == nil {
		temps = []TemporaryConnection{}
	}
	if systems == nil {
		systems = map[int]System{}
	}
	m := &Memory{ansiblexes: ans, tempConnections: temps, systems: systems}
	for _, a := range ans {
		if a.ID >= m.nextAnsiblexID {
			m.nextAnsiblexID = a.ID + 1
		}
	}
	for _, t := range temps {
		if t.ID >= m.nextTempID {
			m.nextTempID = t.ID + 1
		}
	}
	return m
}

// Ansiblexes returns all Ansiblex gates.
func (m *Memory) Ansiblexes(ctx context.Context) ([]Ansiblex, error) {
	return m.ansiblexes, nil
}

// TemporaryConnections returns temporary connections between systems.
func (m *Memory) TemporaryConnections(ctx context.Context) ([]TemporaryConnection, error) {
	return m.tempConnections, nil
}

// Systems returns capital systems information.
func (m *Memory) Systems(ctx context.Context) (map[int]System, error) {
	return m.systems, nil
}

// CreateAnsiblex inserts a new Ansiblex gate.
func (m *Memory) CreateAnsiblex(ctx context.Context, a Ansiblex) (int64, error) {
	a.ID = m.nextAnsiblexID
	m.nextAnsiblexID++
	m.ansiblexes = append(m.ansiblexes, a)
	return a.ID, nil
}

// UpdateAnsiblex updates an existing Ansiblex gate.
func (m *Memory) UpdateAnsiblex(ctx context.Context, a Ansiblex) error {
	for i := range m.ansiblexes {
		if m.ansiblexes[i].ID == a.ID {
			m.ansiblexes[i] = a
			return nil
		}
	}
	return errors.New("ansiblex not found")
}

// DeleteAnsiblex removes an Ansiblex gate.
func (m *Memory) DeleteAnsiblex(ctx context.Context, id int64) error {
	for i := range m.ansiblexes {
		if m.ansiblexes[i].ID == id {
			m.ansiblexes = append(m.ansiblexes[:i], m.ansiblexes[i+1:]...)
			return nil
		}
	}
	return errors.New("ansiblex not found")
}

// CreateTemporaryConnection inserts a new temporary connection.
func (m *Memory) CreateTemporaryConnection(ctx context.Context, c TemporaryConnection) (int64, error) {
	c.ID = m.nextTempID
	m.nextTempID++
	m.tempConnections = append(m.tempConnections, c)
	return c.ID, nil
}

// UpdateTemporaryConnection updates an existing temporary connection.
func (m *Memory) UpdateTemporaryConnection(ctx context.Context, c TemporaryConnection) error {
	for i := range m.tempConnections {
		if m.tempConnections[i].ID == c.ID {
			m.tempConnections[i] = c
			return nil
		}
	}
	return errors.New("temporary connection not found")
}

// DeleteTemporaryConnection removes a temporary connection.
func (m *Memory) DeleteTemporaryConnection(ctx context.Context, id int64) error {
	for i := range m.tempConnections {
		if m.tempConnections[i].ID == id {
			m.tempConnections = append(m.tempConnections[:i], m.tempConnections[i+1:]...)
			return nil
		}
	}
	return errors.New("temporary connection not found")
}
