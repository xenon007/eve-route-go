package graph

import "strings"

// Helper предоставляет вспомогательные функции для работы с графом.
type Helper struct {
	graph Graph
}

// NewHelper создаёт новый экземпляр Helper.
func NewHelper(g Graph) *Helper {
	return &Helper{graph: g}
}

// Graph возвращает исходный граф.
func (h *Helper) Graph() Graph { return h.graph }

// FindSystemByName ищет систему по имени (без учёта регистра).
func (h *Helper) FindSystemByName(name string) *System {
	for _, s := range h.graph.Systems {
		if strings.EqualFold(s.Name, name) {
			sCopy := s
			return &sCopy
		}
	}
	return nil
}

// FindSystem ищет систему по ID.
func (h *Helper) FindSystem(id int) *System {
	for _, s := range h.graph.Systems {
		if s.ID == id {
			sCopy := s
			return &sCopy
		}
	}
	return nil
}

// GetEndSystem возвращает конечную систему из названия Ansiblex.
// Формат названия: "Start » End - ...".
func (h *Helper) GetEndSystem(ansiblexName string) *System {
	parts := strings.Split(ansiblexName, " » ")
	if len(parts) < 2 {
		return nil
	}
	end := parts[1]
	if idx := strings.Index(end, " - "); idx != -1 {
		end = end[:idx]
	}
	return h.FindSystemByName(end)
}
