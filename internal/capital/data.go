package capital

import (
	"embed"
	"encoding/json"
	"log"
)

//go:embed systems.json
var systemsFile embed.FS

// loadSystems загружает данные о системах из встроенного файла.
func loadSystems() map[int]System {
	data, err := systemsFile.ReadFile("systems.json")
	if err != nil {
		log.Fatalf("cannot read systems data: %v", err)
	}
	var list []System
	if err := json.Unmarshal(data, &list); err != nil {
		log.Fatalf("cannot unmarshal systems data: %v", err)
	}
	systems := make(map[int]System)
	for _, s := range list {
		systems[s.ID] = s
	}
	return systems
}

// DefaultSystems возвращает встроенный набор систем для планировщика.
func DefaultSystems() map[int]System {
	return loadSystems()
}
