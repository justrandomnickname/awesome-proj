package aggregates

import (
	"time"
	"awesome-proj/app/domain/entities"
)

// World represents the game world state
type World struct {
	Name      string                         `json:"name"`
	Locations map[string]*entities.Location  `json:"locations"`
	Seed      int64                          `json:"seed"`
}

// NewEmptyWorld creates an empty world without content generation
func NewEmptyWorld(name string, seed int64) *World {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	
	return &World{
		Name:      name,
		Locations: make(map[string]*entities.Location),
		Seed:      seed,
	}
}

// GetLocations returns the world's locations
func (w *World) GetLocations() map[string]*entities.Location {
	return w.Locations
}

// AddLocation adds a location to the world
func (w *World) AddLocation(id string, location *entities.Location) {
	w.Locations[id] = location
}
