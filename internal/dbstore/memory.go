package dbstore

import "context"

// Memory provides an in-memory implementation of Store for tests.
type Memory struct {
	ansiblexes      []Ansiblex
	tempConnections []TemporaryConnection
	systems         map[int]System
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
	return &Memory{ansiblexes: ans, tempConnections: temps, systems: systems}
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
