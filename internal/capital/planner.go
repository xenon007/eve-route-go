package capital

import (
	"container/list"
	"errors"
	"log"
)

// System описывает солнечную систему для капитального маршрута.
type System struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Planner рассчитывает маршрут прыжков капитальных кораблей.
// Используется простой поиск в ширину (BFS).
type Planner struct {
	systems  map[int]System
	graph    map[int][]int
	nameToID map[string]int
}

// NewPlanner создаёт новый планировщик на основе списка систем и графа.
func NewPlanner(systems map[int]System, graph map[int][]int) *Planner {
	nameToID := map[string]int{}
	for id, s := range systems {
		nameToID[s.Name] = id
	}
	return &Planner{systems: systems, graph: graph, nameToID: nameToID}
}

// DefaultSystems содержит минимальный набор систем для примера.
func DefaultSystems() map[int]System {
	return map[int]System{
		1: {ID: 1, Name: "Maila"},
		2: {ID: 2, Name: "Todifrauan"},
	}
}

// DefaultGraph содержит минимальный набор соединений для примера.
func DefaultGraph() map[int][]int {
	return map[int][]int{
		1: {2},
		2: {1},
	}
}

// Plan ищет кратчайший маршрут между системами startName и endName.
func (p *Planner) Plan(startName, endName string) ([]System, error) {
	log.Printf("capital planner: %s -> %s", startName, endName)
	startID, okStart := p.nameToID[startName]
	endID, okEnd := p.nameToID[endName]
	if !okStart || !okEnd {
		return nil, errors.New("system not found")
	}
	if startID == endID {
		return []System{p.systems[startID]}, nil
	}
	visited := map[int]bool{startID: true}
	parent := map[int]int{}
	q := list.New()
	q.PushBack(startID)

	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		s := e.Value.(int)
		for _, n := range p.graph[s] {
			if !visited[n] {
				visited[n] = true
				parent[n] = s
				if n == endID {
					// восстановление пути
					path := []System{p.systems[endID]}
					cur := endID
					for cur != startID {
						cur = parent[cur]
						path = append([]System{p.systems[cur]}, path...)
					}
					return path, nil
				}
				q.PushBack(n)
			}
		}
	}
	return nil, errors.New("route not found")
}
