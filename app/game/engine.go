package game

import (
	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/domain/services"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type GameEngine struct {
	ctx                context.Context
	currentWorld       *aggregates.World
	gameState          *GameState
	saveService        *services.SaveService
	interactionService *services.LocationInteractionService
	isRunning          bool
}

func NewGameEngine() *GameEngine {
	saveService := services.NewSaveService()
	interactionService := services.NewLocationInteractionService()
	world, gameState := initializeGameState(saveService)
	return &GameEngine{
		currentWorld:       world,
		gameState:          gameState,
		saveService:        saveService,
		interactionService: interactionService,
		isRunning:          false,
	}
}

func initializeGameState(saveService *services.SaveService) (*aggregates.World, *GameState) {
	world, gameState, err := tryLoadExistingSave(saveService)
	if err != nil {
		return createNewGame()
	}
	return world, gameState
}

func tryLoadExistingSave(saveService *services.SaveService) (*aggregates.World, *GameState, error) {
	loadedWorld, loadedGameState, err := saveService.LoadFirstAvailableSave()
	if err != nil {
		return nil, nil, err
	}
	gameState, err := deserializeGameState(loadedGameState)
	if err != nil {
		return nil, nil, err
	}
	return loadedWorld, gameState, nil
}

func createNewGame() (*aggregates.World, *GameState) {
	worldService := services.NewWorldGenerationService()
	randomSeed := time.Now().UnixNano()
	world := worldService.GenerateWorld("Default World", randomSeed)
	gameState := NewGameState(randomSeed)
	return world, gameState
}

func deserializeGameState(loadedGameState interface{}) (*GameState, error) {
	jsonBytes, err := json.Marshal(loadedGameState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game state: %v", err)
	}
	var gameState GameState
	if err := json.Unmarshal(jsonBytes, &gameState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %v", err)
	}
	return &gameState, nil
}

func (g *GameEngine) Initialize(ctx context.Context) error {
	g.ctx = ctx
	g.isRunning = true
	return nil
}

func (g *GameEngine) GetCurrentLocationInfo() (*entities.Location, error) {
	if g.gameState == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	currentLocationID := g.gameState.GetCurrentLocationID()
	location := g.currentWorld.Locations[currentLocationID]
	if location == nil {
		return nil, fmt.Errorf("location %s not found", currentLocationID)
	}

	npcs := make([]*entities.NPC, 0)
	for _, npcID := range location.NPCs {
		if npc, exists := g.currentWorld.NPCs[npcID]; exists {
			npcs = append(npcs, npc)
		}
	}

	locationState := g.gameState.GetLocationState(currentLocationID)
	if locationState.FirstVisit && len(locationState.Interactions) == 0 {
		initialState := g.interactionService.GenerateInitialLocationState(
			currentLocationID,
			location.Name,
			len(npcs),
		)
		g.gameState.AddInteractionToCurrentLocation(initialState)
	}

	frontendLocation := *location
	frontendLocation.ToFrontendFormat(npcs, locationState.Interactions)
	return &frontendLocation, nil
}

func (g *GameEngine) IsRunning() bool {
	return g.isRunning
}

func (g *GameEngine) SaveGame(saveName string) error {
	return g.saveService.SaveGame(saveName, g.currentWorld, g.gameState)
}

func (g *GameEngine) LoadGame(filename string) error {
	world, loadedGameState, err := g.saveService.LoadGame(filename)
	if err != nil {
		return err
	}

	g.currentWorld = world
	gameState, err := deserializeGameState(loadedGameState)
	if err != nil {
		return fmt.Errorf("invalid game state format in save file: %v", err)
	}

	g.gameState = gameState
	return nil
}

func (g *GameEngine) GetSavesList() ([]entities.SaveInfo, error) {
	return g.saveService.GetSavesList()
}

func (g *GameEngine) DeleteSave(filename string) error {
	return g.saveService.DeleteSave(filename)
}

func (g *GameEngine) NewGame(seed int64) error {
	worldService := services.NewWorldGenerationService()
	g.currentWorld = worldService.GenerateWorld("Default World", seed)
	g.gameState = NewGameState(seed)
	fmt.Println("[GameEngine] New game started successfully")
	return nil
}

func (g *GameEngine) PerformPlayerAction(actionText string) error {
	if g.gameState == nil {
		return fmt.Errorf("game state not initialized")
	}

	currentLocationID := g.gameState.GetCurrentLocationID()
	location := g.currentWorld.Locations[currentLocationID]
	if location == nil {
		return fmt.Errorf("location %s not found", currentLocationID)
	}

	playerAction := entities.Interaction{
		ID:         fmt.Sprintf("action_%d", time.Now().UnixNano()),
		Type:       entities.InteractionTypePlayerAction,
		Content:    actionText,
		LocationID: currentLocationID,
		Timestamp:  time.Now(),
	}

	g.gameState.AddInteractionToCurrentLocation(playerAction)
	npcCount := len(location.NPCs)
	locationResponse := g.interactionService.GenerateResponseToAction(playerAction, location.Name, npcCount)
	g.gameState.AddInteractionToCurrentLocation(locationResponse)
	return nil
}
