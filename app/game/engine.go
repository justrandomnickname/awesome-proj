package game

import (
	"context"
	"fmt"

	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/domain/services"
)

// GameEngine - центральный игровой движок-оркестратор
type GameEngine struct {
	ctx           context.Context
	currentPlayer *entities.NPC // игрок теперь тоже NPC
	currentWorld  *aggregates.World
	isRunning     bool
}

// NewGameEngine creates a new game engine instance
func NewGameEngine() *GameEngine {
	// Use WorldGenerationService for proper DDD approach
	worldService := services.NewWorldGenerationService()
	world := worldService.GenerateWorld("Default World", 0)
	
	return &GameEngine{
		currentWorld: world,
		isRunning:    false,
	}
}

// Initialize initializes the game engine
func (g *GameEngine) Initialize(ctx context.Context) error {
	g.ctx = ctx
	g.isRunning = true
	
	return nil
}

// CreatePlayer creates a new player (as NPC)
func (g *GameEngine) CreatePlayer(name string) *entities.NPC {
	player := &entities.NPC{
		ID:          fmt.Sprintf("player_%s", name),
		Name:        name,
		Race:        "human", // можно будет выбирать позже
		LocationID:  "start",
		Description: "Это вы - главный герой этой истории",
		Level:       1,
	}
	
	g.currentPlayer = player
	return player
}

// GetCurrentPlayer returns the current player
func (g *GameEngine) GetCurrentPlayer() *entities.NPC {
	return g.currentPlayer
}

// GetCurrentLocation returns player's current location
func (g *GameEngine) GetCurrentLocation() *entities.Location {
	if g.currentPlayer == nil {
		return nil
	}
	
	return g.currentWorld.Locations[g.currentPlayer.LocationID]
}

// MovePlayer moves player to a new location
func (g *GameEngine) MovePlayer(direction string) (*entities.Location, error) {
	currentLocation := g.GetCurrentLocation()
	if currentLocation == nil {
		return nil, fmt.Errorf("player has no current location")
	}
	
	nextLocationID, exists := currentLocation.Exits[direction]
	if !exists {
		return nil, fmt.Errorf("cannot go %s from here", direction)
	}
	
	nextLocation := g.currentWorld.Locations[nextLocationID]
	if nextLocation == nil {
		return nil, fmt.Errorf("location %s does not exist", nextLocationID)
	}
	
	g.currentPlayer.LocationID = nextLocationID
	return nextLocation, nil
}

// GetNPCsInLocation returns all NPCs in the given location
func (g *GameEngine) GetNPCsInLocation(locationID string) []string {
	if location := g.currentWorld.Locations[locationID]; location != nil {
		return location.NPCs
	}
	return []string{}
}

// IsRunning returns whether the game engine is running
func (g *GameEngine) IsRunning() bool {
	return g.isRunning
}
