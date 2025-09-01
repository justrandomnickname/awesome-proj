package aggregates

import (
	"awesome-proj/app/domain/entities"
	"time"
)

type World struct {
	Name      string                        `json:"name"`
	Locations map[string]*entities.Location `json:"locations"`
	NPCs      map[string]*entities.NPC      `json:"npcs"`
	Seed      int64                         `json:"seed"`
}

func NewEmptyWorld(name string, seed int64) *World {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &World{
		Name:      name,
		Locations: make(map[string]*entities.Location),
		NPCs:      make(map[string]*entities.NPC),
		Seed:      seed,
	}
}

func (w *World) GetLocations() map[string]*entities.Location {
	return w.Locations
}

func (w *World) AddLocation(id string, location *entities.Location) {
	w.Locations[id] = location
}

func (w *World) AddNPC(npc *entities.NPC) {
	w.NPCs[npc.ID] = npc
}
