package route

import (
	dbstore "github.com/tkhamez/eve-route-go/internal/dbstore"
	"github.com/tkhamez/eve-route-go/internal/graph"
)

// Ansiblex описывает Ansiblex-ворота.
type Ansiblex = dbstore.Ansiblex

// TemporaryConnection описывает временное соединение между системами.
type TemporaryConnection = dbstore.TemporaryConnection

// ConnectedSystems — пара систем, связь между которыми удалена пользователем.
type ConnectedSystems struct {
	System1 string
	System2 string
}

// Waypoint описывает один шаг маршрута.
type Waypoint struct {
	SystemID       int
	SystemName     string
	TargetSystem   *string
	Wormhole       bool
	SystemSecurity float64
	ConnectionType *WaypointType
	AnsiblexID     *int64
	AnsiblexName   *string
	RegionName     string
}

// WaypointType тип соединения.
type WaypointType string

const (
	TypeStargate  WaypointType = "Stargate"
	TypeAnsiblex  WaypointType = "Ansiblex"
	TypeTemporary WaypointType = "Temporary"
)

// Helper alias для системного типа из пакета graph.
type GraphSystem = graph.System
