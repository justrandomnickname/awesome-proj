package game

import (
	"context"
	"fmt"

	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/services"
)

// GameEngine - центральный игровой движок-оркестратор
type GameEngine struct {
	ctx          context.Context
	currentWorld *aggregates.World
	gameState    *GameState
	isRunning    bool
}

// NewGameEngine creates a new game engine instance
func NewGameEngine() *GameEngine {
	// Use WorldGenerationService for proper DDD approach
	worldService := services.NewWorldGenerationService()
	world := worldService.GenerateWorld("Default World", 0)
	
	return &GameEngine{
		currentWorld: world,
		gameState:    NewGameState(),
		isRunning:    false,
	}
}

// Initialize initializes the game engine
func (g *GameEngine) Initialize(ctx context.Context) error {
	g.ctx = ctx
	g.isRunning = true
	
	return nil
}

// GetCurrentLocationInfo returns current location info for frontend
func (g *GameEngine) GetCurrentLocationInfo() (*LocationInfo, error) {
	if g.gameState == nil {
		return nil, fmt.Errorf("game state not initialized")
	}
	
	currentLocationID := g.gameState.GetCurrentLocationID()
	location := g.currentWorld.Locations[currentLocationID]
	
	if location == nil {
		return nil, fmt.Errorf("location %s not found", currentLocationID)
	}
	
	// Get NPCs in this location
	npcInfos := make([]NPCInfo, 0)
	for _, npcID := range location.NPCs {
		if npc, exists := g.currentWorld.NPCs[npcID]; exists {
			npcInfo := NPCInfo{
				ID:          npc.ID,
				Name:        npc.Name,
				Race:        npc.Race,
				Description: npc.Description,
			}
			npcInfos = append(npcInfos, npcInfo)
		}
	}
	
	return &LocationInfo{
		ID:          location.ID,
		Name:        location.Name,
		Description: location.Description,
		NPCs:        npcInfos,
	}, nil
}

// IsRunning returns whether the game engine is running
func (g *GameEngine) IsRunning() bool {
	return g.isRunning
}
