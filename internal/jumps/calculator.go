package jumps

import (
	"container/list"
	"errors"
	"github.com/tkhamez/eve-route-go/internal/graph"
	"log"
)

// Calculator рассчитывает количество прыжков между системами.
type Calculator struct {
	helper *graph.Helper
	adj    map[int][]int
}

// NewCalculator создаёт новый калькулятор на основе графа.
func NewCalculator(g graph.Graph) *Calculator {
	h := graph.NewHelper(g)
	adj := map[int][]int{}
	for _, c := range g.Connections {
		adj[c[0]] = append(adj[c[0]], c[1])
		adj[c[1]] = append(adj[c[1]], c[0])
	}
	return &Calculator{helper: h, adj: adj}
}

// Between возвращает минимальное число прыжков от from до to.
func (c *Calculator) Between(from, to string) (int, error) {
	log.Printf("jumps: %s -> %s", from, to)
	start := c.helper.FindSystemByName(from)
	end := c.helper.FindSystemByName(to)
	if start == nil || end == nil {
		return 0, errors.New("system not found")
	}
	if start.ID == end.ID {
		return 0, nil
	}
	visited := map[int]bool{start.ID: true}
	dist := map[int]int{start.ID: 0}
	q := list.New()
	q.PushBack(start.ID)
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		cur := e.Value.(int)
		if cur == end.ID {
			return dist[cur], nil
		}
		for _, n := range c.adj[cur] {
			if !visited[n] {
				visited[n] = true
				dist[n] = dist[cur] + 1
				q.PushBack(n)
			}
		}
	}
	return 0, errors.New("route not found")
}
