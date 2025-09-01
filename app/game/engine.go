package game

import (
	"context"
	"fmt"
	"time"

	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/domain/services"
)

// GameEngine - центральный игровой движок-оркестратор
type GameEngine struct {
	ctx                 context.Context
	currentWorld        *aggregates.World
	gameState          *GameState
	saveService        *services.SaveService
	interactionService *services.LocationInteractionService
	isRunning          bool
}

// NewGameEngine creates a new game engine instance
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

// initializeGameState попытка загрузить сохранение или создать новую игру
func initializeGameState(saveService *services.SaveService) (*aggregates.World, *GameState) {
	world, gameState, err := tryLoadExistingSave(saveService)
	if err != nil {
		return createNewGame()
	}
	return world, gameState
}

// tryLoadExistingSave попытка загрузить первый доступный сейв
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

// createNewGame создаёт новую игру с случайным сидом
func createNewGame() (*aggregates.World, *GameState) {
	worldService := services.NewWorldGenerationService()
	randomSeed := time.Now().UnixNano()
	world := worldService.GenerateWorld("Default World", randomSeed)
	gameState := NewGameState(randomSeed)
	return world, gameState
}

// deserializeGameState десериализует GameState из interface{}
func deserializeGameState(loadedGameState interface{}) (*GameState, error) {
	gsMap, ok := loadedGameState.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid game state format")
	}
	
	worldSeed := extractWorldSeed(gsMap)
	locationStates := extractLocationStates(gsMap)
	currentLocationID := extractCurrentLocationID(gsMap)
	
	return &GameState{
		CurrentLocationID: currentLocationID,
		WorldSeed:         worldSeed,
		LocationStates:    locationStates,
	}, nil
}

// extractWorldSeed извлекает seed мира из данных сохранения
func extractWorldSeed(gsMap map[string]interface{}) int64 {
	if seedFloat, ok := gsMap["world_seed"].(float64); ok {
		return int64(seedFloat)
	}
	return 0
}

// extractCurrentLocationID извлекает текущую локацию из данных сохранения
func extractCurrentLocationID(gsMap map[string]interface{}) string {
	if locationID, ok := gsMap["current_location_id"].(string); ok {
		return locationID
	}
	return ""
}

// extractLocationStates извлекает состояния локаций из данных сохранения
func extractLocationStates(gsMap map[string]interface{}) map[string]*entities.LocationState {
	locationStates := make(map[string]*entities.LocationState)
	
	locationStatesData, ok := gsMap["location_states"].(map[string]interface{})
	if !ok {
		return locationStates
	}
	
	for locationID, stateData := range locationStatesData {
		locationState := parseLocationState(locationID, stateData)
		if locationState != nil {
			locationStates[locationID] = locationState
		}
	}
	
	return locationStates
}

// parseLocationState парсит отдельное состояние локации
func parseLocationState(locationID string, stateData interface{}) *entities.LocationState {
	stateMap, ok := stateData.(map[string]interface{})
	if !ok {
		return nil
	}
	
	locationState := &entities.LocationState{
		LocationID:   locationID,
		Interactions: make([]entities.Interaction, 0),
		FirstVisit:   true,
	}
	
	if firstVisit, ok := stateMap["first_visit"].(bool); ok {
		locationState.FirstVisit = firstVisit
	}
	
	locationState.Interactions = parseInteractions(stateMap)
	
	return locationState
}

// parseInteractions парсит массив взаимодействий
func parseInteractions(stateMap map[string]interface{}) []entities.Interaction {
	var interactions []entities.Interaction
	
	interactionsData, ok := stateMap["interactions"].([]interface{})
	if !ok {
		return interactions
	}
	
	for _, interactionData := range interactionsData {
		interaction := parseInteraction(interactionData)
		if interaction != nil {
			interactions = append(interactions, *interaction)
		}
	}
	
	return interactions
}

// parseInteraction парсит отдельное взаимодействие
func parseInteraction(interactionData interface{}) *entities.Interaction {
	interactionMap, ok := interactionData.(map[string]interface{})
	if !ok {
		return nil
	}
	
	interaction := entities.Interaction{
		ID:         interactionMap["id"].(string),
		Type:       entities.InteractionType(interactionMap["type"].(string)),
		Content:    interactionMap["content"].(string),
		LocationID: interactionMap["location_id"].(string),
	}
	
	if timestampStr, ok := interactionMap["timestamp"].(string); ok {
		if timestamp, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			interaction.Timestamp = timestamp
		}
	}
	
	return &interaction
}

// Initialize initializes the game engine
func (g *GameEngine) Initialize(ctx context.Context) error {
	g.ctx = ctx
	g.isRunning = true
	
	return nil
}

// GetCurrentLocationInfo returns current location info for frontend
func (g *GameEngine) GetCurrentLocationInfo() (*entities.Location, error) {
	if g.gameState == nil {
		return nil, fmt.Errorf("game state not initialized")
	}
	
	currentLocationID := g.gameState.GetCurrentLocationID()
	location := g.currentWorld.Locations[currentLocationID]
	
	if location == nil {
		return nil, fmt.Errorf("location %s not found", currentLocationID)
	}
	
	// Get NPCs in this location
	npcs := make([]*entities.NPC, 0)
	for _, npcID := range location.NPCs {
		if npc, exists := g.currentWorld.NPCs[npcID]; exists {
			npcs = append(npcs, npc)
		}
	}
	
	// Get or create location state
	locationState := g.gameState.GetLocationState(currentLocationID)
	
	// Если это первое посещение - генерируем начальное состояние локации
	if locationState.FirstVisit && len(locationState.Interactions) == 0 {
		initialState := g.interactionService.GenerateInitialLocationState(
			currentLocationID, 
			location.Name, 
			len(npcs),
		)
		g.gameState.AddInteractionToCurrentLocation(initialState)
	}
	
	// Создаем копию location для frontend и заполняем детальными данными
	frontendLocation := *location
	frontendLocation.ToFrontendFormat(npcs, locationState.Interactions)
	
	return &frontendLocation, nil
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
func (g *GameEngine) GetSavesList() ([]entities.SaveInfo, error) {
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

// PerformPlayerAction handles player action and generates response
func (g *GameEngine) PerformPlayerAction(actionText string) error {
	if g.gameState == nil {
		return fmt.Errorf("game state not initialized")
	}
	
	currentLocationID := g.gameState.GetCurrentLocationID()
	location := g.currentWorld.Locations[currentLocationID]
	
	if location == nil {
		return fmt.Errorf("location %s not found", currentLocationID)
	}
	
	// Создаем действие игрока
	playerAction := entities.Interaction{
		ID:         fmt.Sprintf("action_%d", time.Now().UnixNano()),
		Type:       entities.InteractionTypePlayerAction,
		Content:    actionText,
		LocationID: currentLocationID,
		Timestamp:  time.Now(),
	}
	
	// Добавляем действие в состояние локации
	g.gameState.AddInteractionToCurrentLocation(playerAction)
	
	// Генерируем ответ от локации
	npcCount := len(location.NPCs)
	locationResponse := g.interactionService.GenerateResponseToAction(playerAction, location.Name, npcCount)
	
	// Добавляем ответ в состояние локации
	g.gameState.AddInteractionToCurrentLocation(locationResponse)
	
	return nil
}
