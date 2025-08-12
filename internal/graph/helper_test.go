package graph

import (
	"log"
	"reflect"
	"testing"
)

// TestNewHelperAndGraph проверяет создание Helper и возврат графа.
func TestNewHelperAndGraph(t *testing.T) {
	g := DefaultGraph()
	h := NewHelper(g)
	log.Printf("создан helper: %+v", h)

	if !reflect.DeepEqual(h.Graph(), g) {
		t.Errorf("ожидался граф %+v, получено %+v", g, h.Graph())
	}
}

// TestFindSystemByName проверяет поиск системы по имени.
func TestFindSystemByName(t *testing.T) {
	h := NewHelper(DefaultGraph())

	s := h.FindSystemByName("beta")
	log.Printf("найдена система: %+v", s)
	if s == nil || s.Name != "Beta" {
		t.Fatalf("ожидалась система Beta, получено %#v", s)
	}
	if res := h.FindSystemByName("unknown"); res != nil {
		t.Errorf("ожидалось nil, получено %+v", res)
	}
}

// TestFindSystem проверяет поиск системы по ID.
func TestFindSystem(t *testing.T) {
	h := NewHelper(DefaultGraph())
	s := h.FindSystem(2)
	log.Printf("найдена система по ID: %+v", s)
	if s == nil || s.Name != "Beta" {
		t.Fatalf("ожидалась система Beta, получено %#v", s)
	}
	if res := h.FindSystem(999); res != nil {
		t.Errorf("ожидалось nil, получено %+v", res)
	}
}

// TestGetEndSystem проверяет извлечение конечной системы из названия Ansiblex.
func TestGetEndSystem(t *testing.T) {
	h := NewHelper(DefaultGraph())
	ansiblex := "Alpha » Beta - jump bridge"
	end := h.GetEndSystem(ansiblex)
	log.Printf("конечная система: %+v", end)
	if end == nil || end.Name != "Beta" {
		t.Fatalf("ожидалась Beta, получено %#v", end)
	}
	if res := h.GetEndSystem("invalid format"); res != nil {
		t.Errorf("ожидалось nil, получено %+v", res)
	}
	if res := h.GetEndSystem("Alpha » Unknown - test"); res != nil {
		t.Errorf("ожидалось nil при отсутствии системы, получено %+v", res)
	}
}
