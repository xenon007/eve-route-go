package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	routepkg "github.com/tkhamez/eve-route-go/internal/route"
)

// NewRouteHandler возвращает HTTP-обработчик, строящий маршрут между системами.
func NewRouteHandler(r *routepkg.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		from := vars["from"]
		to := vars["to"]
		if from == "" || to == "" {
			http.Error(w, "missing from or to", http.StatusBadRequest)
			return
		}
		paths := r.Find(from, to)
		if len(paths) == 0 {
			http.NotFound(w, req)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"routes": paths})
	}
}
