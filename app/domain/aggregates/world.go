package aggregates

import (
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/domain/entities/npc"
	"time"
)

type World struct {
	Name      string                        `json:"name"`
	Locations map[string]*entities.Location `json:"locations"`
	NPCs      map[string]*npc.NPC           `json:"npcs"`
	Hierarchy *entities.LocationHierarchy   `json:"hierarchy,omitempty"`
	Seed      int64                         `json:"seed"`
}

func NewEmptyWorld(name string, seed int64) *World {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &World{
		Name:      name,
		Locations: make(map[string]*entities.Location),
		NPCs:      make(map[string]*npc.NPC),
		Hierarchy: entities.NewLocationHierarchy(),
		Seed:      seed,
	}
}

func (w *World) GetLocations() map[string]*entities.Location {
	return w.Locations
}

func (w *World) AddLocation(id string, location *entities.Location) {
	w.Locations[id] = location
}

func (w *World) AddNPC(npcEntity *npc.NPC) {
	w.NPCs[npcEntity.ID] = npcEntity
}

func (w *World) SetHierarchy(hierarchy *entities.LocationHierarchy) {
	w.Hierarchy = hierarchy
}

func (w *World) GetHierarchy() *entities.LocationHierarchy {
	return w.Hierarchy
}
