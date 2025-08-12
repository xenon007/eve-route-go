package route

// Node представляет узел графа.
type Node struct {
	Value       GraphSystem
	connections []Connection
}

// Connection соединение между узлами.
type Connection struct {
	Node *Node
	Type WaypointType
}

// Connect соединяет узлы в обе стороны, если они ещё не соединены.
func (n *Node) Connect(other *Node, t WaypointType) {
	if n == other {
		return
	}
	if !n.isConnected(other, t) {
		n.connections = append(n.connections, Connection{Node: other, Type: t})
	}
	if !other.isConnected(n, t) {
		other.connections = append(other.connections, Connection{Node: n, Type: t})
	}
}

func (n *Node) isConnected(other *Node, t WaypointType) bool {
	for _, c := range n.connections {
		if c.Node == other && c.Type == t {
			return true
		}
	}
	return false
}

// Connections возвращает все соединения узла.
func (n *Node) Connections() []Connection { return n.connections }

type path struct {
	waypoints []Waypoint
}

// numberOfAnsiblexes возвращает количество переходов через Ansiblex.
func (p path) numberOfAnsiblexes() int {
	count := 0
	for _, w := range p.waypoints {
		if w.ConnectionType != nil && *w.ConnectionType == TypeAnsiblex {
			count++
		}
	}
	return count
}

// numberOfTemporary возвращает количество переходов через временные соединения.
func (p path) numberOfTemporary() int {
	count := 0
	for _, w := range p.waypoints {
		if w.ConnectionType != nil && *w.ConnectionType == TypeTemporary {
			count++
		}
	}
	return count
}
