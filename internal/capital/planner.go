package capital

import (
	"container/list"
	"errors"
	"log"
	"math"
)

// lyInMeters — количество метров в одном световом годе.
const lyInMeters = 9.4607e15

// System описывает солнечную систему для капитального маршрута.
type System struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

// Planner рассчитывает маршрут прыжков капитальных кораблей.
// Для поиска используется BFS, соседи вычисляются по радиусу прыжка.
type Planner struct {
	systems   map[int]System
	nameToID  map[string]int
	jumpRange float64 // в световых годах
}

// NewPlanner создаёт новый планировщик на основе списка систем и радиуса прыжка.
func NewPlanner(systems map[int]System, jumpRange float64) *Planner {
	nameToID := map[string]int{}
	for id, s := range systems {
		nameToID[s.Name] = id
	}
	return &Planner{systems: systems, nameToID: nameToID, jumpRange: jumpRange}
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
		cur := e.Value.(int)
		if cur == endID {
			break
		}
		for _, n := range p.neighbors(cur) {
			if !visited[n] {
				visited[n] = true
				parent[n] = cur
				q.PushBack(n)
			}
		}
	}

	if !visited[endID] {
		return nil, errors.New("route not found")
	}

	// восстановление пути
	path := []System{p.systems[endID]}
	cur := endID
	for cur != startID {
		cur = parent[cur]
		path = append([]System{p.systems[cur]}, path...)
	}
	return path, nil
}

// neighbors возвращает список систем, достижимых из указанной за один прыжок.
func (p *Planner) neighbors(id int) []int {
	var res []int
	cur := p.systems[id]
	for oid, s := range p.systems {
		if oid == id {
			continue
		}
		if distance(cur, s) <= p.jumpRange {
			res = append(res, oid)
		}
	}
	return res
}

// distance вычисляет расстояние между двумя системами в световых годах.
func distance(a, b System) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx+dy*dy+dz*dz) / lyInMeters
}
