package db

import (
	"context"

	"github.com/tkhamez/eve-route-go/internal/graph"
)

// Memory provides an in-memory implementation of Store for tests.
type Memory struct {
	ansiblexes      []Ansiblex
	tempConnections []TemporaryConnection
	systems         map[int]System
	graph           graph.Graph
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

// Graph returns the stored universe graph.
func (m *Memory) Graph(ctx context.Context) (graph.Graph, error) {
	return m.graph, nil
}

// SaveGraph stores the universe graph in memory.
func (m *Memory) SaveGraph(ctx context.Context, g graph.Graph) error {
	m.graph = g
	return nil
}
