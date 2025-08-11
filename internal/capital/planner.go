package capital

import (
	"container/list"
	"errors"
	"log"
)

// Planner рассчитывает маршрут прыжков капитальных кораблей
// на основе графа систем EVE Online.
type Planner struct {
	graph map[string][]string
}

// NewPlanner создаёт новый планировщик.
func NewPlanner(graph map[string][]string) *Planner {
	return &Planner{graph: graph}
}

// DefaultGraph содержит минимальный набор данных для примера.
func DefaultGraph() map[string][]string {
	return map[string][]string{
		"Maila":      {"Todifrauan"},
		"Todifrauan": {"Maila"},
	}
}

// Plan ищет кратчайший маршрут между системами start и end.
// Используется простой BFS без учёта реальной дальности прыжка.
func (p *Planner) Plan(start, end string) ([]string, error) {
	log.Printf("capital planner: %s -> %s", start, end)
	if start == end {
		return []string{start}, nil
	}
	visited := map[string]bool{start: true}
	parent := map[string]string{}
	q := list.New()
	q.PushBack(start)

	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		s := e.Value.(string)
		for _, n := range p.graph[s] {
			if !visited[n] {
				visited[n] = true
				parent[n] = s
				if n == end {
					// восстановление пути
					path := []string{end}
					for cur := end; cur != start; {
						cur = parent[cur]
						path = append([]string{cur}, path...)
					}
					return path, nil
				}
				q.PushBack(n)
			}
		}
	}
	return nil, errors.New("route not found")
}
