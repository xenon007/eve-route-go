package route

import (
	"context"
	"log"
	"sort"

	"github.com/tkhamez/eve-route-go/internal/db"
	"github.com/tkhamez/eve-route-go/internal/graph"
)

// Route ищет пути между системами на основании графа.
type Route struct {
	graphHelper        *graph.Helper
	avoidedSystems     map[int]bool
	removedConnections []ConnectedSystems

	allSystems              map[int]GraphSystem
	allNodes                map[int]*Node
	allAnsiblexes           map[int]Ansiblex
	allTemporaryConnections map[int]TemporaryConnection
}

// NewRoute создаёт новый экземпляр маршрутизатора и загружает данные из хранилища.
func NewRoute(store db.Store, avoided map[int]bool, removed []ConnectedSystems) (*Route, error) {
	g := graph.DefaultGraph()
	helper := graph.NewHelper(g)
	r := &Route{
		graphHelper:             helper,
		avoidedSystems:          avoided,
		removedConnections:      removed,
		allSystems:              map[int]GraphSystem{},
		allNodes:                map[int]*Node{},
		allAnsiblexes:           map[int]Ansiblex{},
		allTemporaryConnections: map[int]TemporaryConnection{},
	}
	r.buildNodes()
	ansiblexes, err := store.Ansiblexes(context.Background())
	if err != nil {
		return nil, err
	}
	tempConnections, err := store.TemporaryConnections(context.Background())
	if err != nil {
		return nil, err
	}
	r.addGates(ansiblexes)
	r.addTempConnections(tempConnections)
	return r, nil
}

// Find ищет пути от from до to. Возвращает список маршрутов с набором точек.
func (r *Route) Find(from, to string) [][]Waypoint {
	log.Printf("route planner: %s -> %s", from, to)
	startSystem := r.graphHelper.FindSystemByName(from)
	endSystem := r.graphHelper.FindSystemByName(to)
	if startSystem == nil || endSystem == nil {
		return [][]Waypoint{}
	}
	startNode := r.allNodes[startSystem.ID]
	if startNode == nil {
		return [][]Waypoint{}
	}
	connections := r.search(*endSystem, startNode)
	var paths []path
	for _, c := range connections {
		wp := r.buildWaypoints(c)
		paths = append(paths, path{waypoints: wp})
	}
	sort.Slice(paths, func(i, j int) bool {
		ai := paths[i].numberOfAnsiblexes()
		aj := paths[j].numberOfAnsiblexes()
		if ai == aj {
			return paths[i].numberOfTemporary() < paths[j].numberOfTemporary()
		}
		return ai < aj
	})
	var result [][]Waypoint
	for _, p := range paths {
		result = append(result, p.waypoints)
	}
	return result
}

// buildNodes создаёт узлы и соединяет их в соответствии с графом.
func (r *Route) buildNodes() {
	g := r.graphHelper.Graph()
	for _, s := range g.Systems {
		r.allSystems[s.ID] = s
	}
	for _, c := range g.Connections {
		src := r.getNode(c[0])
		dst := r.getNode(c[1])
		if src != nil && dst != nil && !r.isRemoved(src.Value.Name, dst.Value.Name) {
			src.Connect(dst, TypeStargate)
		}
	}
}

func (r *Route) addGates(ansiblexes []Ansiblex) {
	for _, gate := range ansiblexes {
		r.allAnsiblexes[gate.SolarSystemID] = gate
		end := r.graphHelper.GetEndSystem(gate.Name)
		if end == nil {
			continue
		}
		startNode := r.getNode(gate.SolarSystemID)
		endNode := r.getNode(end.ID)
		if startNode != nil && endNode != nil && !r.isRemoved(startNode.Value.Name, endNode.Value.Name) {
			startNode.Connect(endNode, TypeAnsiblex)
		}
	}
}

func (r *Route) addTempConnections(conns []TemporaryConnection) {
	for _, c := range conns {
		r.allTemporaryConnections[c.System1ID] = c
		r.allTemporaryConnections[c.System2ID] = c
		n1 := r.getNode(c.System1ID)
		n2 := r.getNode(c.System2ID)
		if n1 != nil && n2 != nil && !r.isRemoved(n1.Value.Name, n2.Value.Name) {
			n1.Connect(n2, TypeTemporary)
		}
	}
}

func (r *Route) isRemoved(startName, endName string) bool {
	for _, rc := range r.removedConnections {
		if (rc.System1 == startName && rc.System2 == endName) || (rc.System1 == endName && rc.System2 == startName) {
			return true
		}
	}
	return false
}

func (r *Route) getNode(systemID int) *Node {
	if r.avoidedSystems != nil && r.avoidedSystems[systemID] {
		return nil
	}
	if n, ok := r.allNodes[systemID]; ok {
		return n
	}
	if s, ok := r.allSystems[systemID]; ok {
		node := &Node{Value: s}
		r.allNodes[systemID] = node
		return node
	}
	return nil
}

func (r *Route) search(goal GraphSystem, start *Node) [][]Connection {
	type connPath []Connection
	var found []connPath
	queue := []connPath{{{Node: start}}}
	visited := map[*Node]bool{}
	lastLen := 0
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		current := p[len(p)-1]
		if lastLen > 0 && len(p) > lastLen {
			continue
		}
		if current.Node.Value.ID == goal.ID {
			if lastLen == 0 || len(p) < lastLen {
				found = []connPath{p}
				lastLen = len(p)
			} else if len(p) == lastLen {
				found = append(found, p)
			}
			continue
		}
		if visited[current.Node] {
			continue
		}
		visited[current.Node] = true
		for _, c := range current.Node.Connections() {
			newPath := append(connPath{}, p...)
			newPath = append(newPath, c)
			queue = append(queue, newPath)
		}
	}
	var result [][]Connection
	for _, p := range found {
		result = append(result, []Connection(p))
	}
	return result
}

func (r *Route) buildWaypoints(path []Connection) []Waypoint {
	var waypoints []Waypoint
	for i := len(path) - 1; i >= 0; i-- {
		conn := path[i]
		system := conn.Node.Value
		var prevSystem *GraphSystem
		if i < len(path)-1 {
			prevSystem = &path[i+1].Node.Value
		}
		var ansiblexID *int64
		var ansiblexName *string
		if i < len(path)-1 && path[i+1].Type == TypeAnsiblex {
			if gate, ok := r.allAnsiblexes[system.ID]; ok {
				ansiblexID = &gate.ID
				ansiblexName = &gate.Name
			}
		}
		var prevName *string
		if prevSystem != nil {
			name := prevSystem.Name
			prevName = &name
		}
		w := Waypoint{
			SystemID:       system.ID,
			SystemName:     system.Name,
			TargetSystem:   prevName,
			Wormhole:       system.ID >= 31000000 && system.ID <= 32000000,
			SystemSecurity: system.Security,
			RegionName:     r.graphHelper.Graph().Regions[system.RegionID],
		}
		if i < len(path)-1 {
			t := path[i+1].Type
			w.ConnectionType = &t
			w.AnsiblexID = ansiblexID
			w.AnsiblexName = ansiblexName
		}
		waypoints = append(waypoints, w)
	}
	// reverse
	for i, j := 0, len(waypoints)-1; i < j; i, j = i+1, j-1 {
		waypoints[i], waypoints[j] = waypoints[j], waypoints[i]
	}
	return waypoints
}
