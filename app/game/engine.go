package game

import (
	"context"
	"fmt"
	"time"

	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/services"
)

// GameEngine - центральный игровой движок-оркестратор
type GameEngine struct {
	ctx          context.Context
	currentWorld *aggregates.World
	gameState    *GameState
	saveService  *services.SaveService
	isRunning    bool
}

// NewGameEngine creates a new game engine instance
func NewGameEngine() *GameEngine {
	saveService := services.NewSaveService()
	
	var world *aggregates.World
	var gameState *GameState
	
	// Пытаемся загрузить первый доступный сейв
	loadedWorld, loadedGameState, err := saveService.LoadFirstAvailableSave()
	if err != nil {
		// Если нет сейвов - создаем новую игру
		worldService := services.NewWorldGenerationService()
		randomSeed := time.Now().UnixNano()
		world = worldService.GenerateWorld("Default World", randomSeed)
		gameState = NewGameState(randomSeed)
	} else {
		world = loadedWorld
		// Десериализуем GameState из interface{}
		if gsMap, ok := loadedGameState.(map[string]interface{}); ok {
			worldSeed := int64(0)
			if seedFloat, ok := gsMap["world_seed"].(float64); ok {
				worldSeed = int64(seedFloat)
			}
			gameState = &GameState{
				CurrentLocationID: gsMap["current_location_id"].(string),
				WorldSeed:         worldSeed,
			}
		} else {			// Если не можем десериализовать - создаем новое состояние
			randomSeed := time.Now().UnixNano()
			gameState = NewGameState(randomSeed)
		}
	}
	
	return &GameEngine{
		currentWorld: world,
		gameState:    gameState,
		saveService:  saveService,
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

// SaveGame saves the current game state
func (g *GameEngine) SaveGame(saveName string) error {
	return g.saveService.SaveGame(saveName, g.currentWorld, g.gameState)
}

// LoadGame loads a game from save file
func (g *GameEngine) LoadGame(filename string) error {
	world, loadedGameState, err := g.saveService.LoadGame(filename)
	if err != nil {
		return err
	}
	
	g.currentWorld = world
	
	// Десериализуем GameState из interface{}
	if gsMap, ok := loadedGameState.(map[string]interface{}); ok {
		g.gameState = &GameState{
			CurrentLocationID: gsMap["current_location_id"].(string),
		}
	} else {
		return fmt.Errorf("invalid game state format in save file")
	}
	
	return nil
}

// GetSavesList returns list of available saves
func (g *GameEngine) GetSavesList() ([]services.SaveInfo, error) {
	return g.saveService.GetSavesList()
}

// DeleteSave deletes a save file
func (g *GameEngine) DeleteSave(filename string) error {
	return g.saveService.DeleteSave(filename)
}

// NewGame starts a new game with specified seed
func (g *GameEngine) NewGame(seed int64) error {
	fmt.Printf("[GameEngine] Starting new game with seed %d...\n", seed)
	
	// Создаём новый мир с переданным сидом
	worldService := services.NewWorldGenerationService()
	g.currentWorld = worldService.GenerateWorld("Default World", seed)
	g.gameState = NewGameState(seed)
	fmt.Printf("[GameEngine] Created new world with seed %d\n", seed)
	
	// Не удаляем сохранения - пользователь может хотеть сохранить новую игру
	fmt.Println("[GameEngine] New game started successfully")
	return nil
}
