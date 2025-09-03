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
		// Если нет сохранений или они битые, создаем новую игру
		// НО только если это действительно отсутствие сохранений
		if saveService.HasAnySaves() {
			// Если сохранения есть, но они битые - это критическая ошибка
			panic(fmt.Sprintf("Failed to load existing save: %v. This indicates corrupted save data.", err))
		}
		// Если сохранений нет вообще - создаем новую игру
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

	// Инициализируем начальный point
	err := g.initializeStartingPoint()
	if err != nil {
		return fmt.Errorf("failed to initialize starting point: %v", err)
	}

	return nil
}

func (g *GameEngine) initializeStartingPoint() error {
	hierarchy := g.currentWorld.GetHierarchy()
	if hierarchy == nil {
		return fmt.Errorf("hierarchy not available")
	}

	// Пытаемся найти указанную точку
	currentPointID := g.gameState.GetCurrentPointID()
	point := hierarchy.FindPoint(currentPointID)

	if point != nil {
		return nil // Точка найдена, все в порядке
	}

	// Если указанная точка не найдена, ищем первую entry point
	for _, cluster := range hierarchy.Clusters {
		for _, subCluster := range cluster.SubClusters {
			for _, entryPointID := range subCluster.EntryPoints {
				if entryPoint := hierarchy.FindPoint(entryPointID); entryPoint != nil {
					g.gameState.SetCurrentPointID(entryPointID)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no entry points found in hierarchy")
}

func (g *GameEngine) PerformPlayerInteraction(actionText string, interactionType entities.InteractionType, additionalContent string) error {
	if g.gameState == nil {
		return fmt.Errorf("game state not initialized")
	}

	currentPoint, err := g.GetCurrentPoint()
	if currentPoint == nil {
		return fmt.Errorf("current point not found: %v", err)
	}

	playerAction := entities.Interaction{
		ID:                fmt.Sprintf("action_%d", time.Now().UnixNano()),
		Type:              interactionType,
		Content:           actionText,
		LocationID:        currentPoint.ID, // Используем Point ID вместо Location ID
		AdditionalContent: additionalContent,
		Timestamp:         time.Now(),
	}

	g.gameState.AddInteractionToCurrentPoint(playerAction)
	npcCount := len(currentPoint.NPCs)
	locationResponse := g.interactionService.GenerateResponseToAction(playerAction, currentPoint.Name, npcCount)
	g.gameState.AddInteractionToCurrentPoint(locationResponse)
	return nil
}

func (g *GameEngine) PerformPlayerAction(actionText string) error {
	return g.PerformPlayerInteraction(actionText, entities.InteractionTypePlayerAction, "")
}

func (g *GameEngine) PerformPlayerMovement(actionText string, targetPoint *entities.Point) error {
	additionalContent := targetPoint.Description
	return g.PerformPlayerInteraction(actionText, entities.InteractionTypePlayerMovement, additionalContent)
}

func (g *GameEngine) GetLocationHierarchy() (*entities.LocationHierarchy, error) {
	if g.currentWorld == nil {
		return nil, fmt.Errorf("world not initialized")
	}
	return g.currentWorld.GetHierarchy(), nil
}
func (g *GameEngine) GetCurrentPoint() (*entities.Point, error) {
	if g.gameState == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	hierarchy := g.currentWorld.GetHierarchy()
	if hierarchy == nil {
		return nil, fmt.Errorf("hierarchy not available")
	}

	currentPointID := g.gameState.GetCurrentPointID()
	if currentPointID == "" {
		// Если CurrentPointID не установлен, инициализируем его
		err := g.initializeStartingPoint()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize starting point: %v", err)
		}
		currentPointID = g.gameState.GetCurrentPointID()
	}

	point := hierarchy.FindPoint(currentPointID)
	if point == nil {
		return nil, fmt.Errorf("current point %s not found in hierarchy", currentPointID)
	}

	return point, nil
}

func (g *GameEngine) MoveToPoint(pointID string) error {
	if g.gameState == nil {
		return fmt.Errorf("game state not initialized")
	}

	hierarchy := g.currentWorld.GetHierarchy()
	if hierarchy == nil {
		return fmt.Errorf("hierarchy not available")
	}

	targetPoint := hierarchy.FindPoint(pointID)
	if targetPoint == nil {
		return fmt.Errorf("point %s not found", pointID)
	}

	currentPoint, err := g.GetCurrentPoint()
	if err != nil {
		return fmt.Errorf("failed to get current point: %v", err)
	}

	canMove := false

	//TODO: переделать наа slices.Contains
	for _, connectionID := range currentPoint.Connections {
		if connectionID == pointID {
			canMove = true
			break
		}
	}

	if !canMove {
		return fmt.Errorf("cannot move to point %s from current location", pointID)
	}

	g.PerformPlayerMovement(fmt.Sprintf("Персонаж перешел из %s в %s", currentPoint.Name, targetPoint.Name), targetPoint)
	g.gameState.SetCurrentPointID(pointID)
	return nil
}
func (g *GameEngine) GetAvailableConnections() ([]*entities.Point, error) {
	currentPoint, err := g.GetCurrentPoint()
	if err != nil {
		return nil, err
	}

	hierarchy := g.currentWorld.GetHierarchy()
	connections := make([]*entities.Point, 0, len(currentPoint.Connections))

	for _, connectionID := range currentPoint.Connections {
		point := hierarchy.FindPoint(connectionID)
		if point != nil {
			connections = append(connections, point)
		}
	}

	return connections, nil
}

func (g *GameEngine) GetAvailableConnectionsInfo() ([]*entities.ConnectionInfo, error) {
	currentPoint, err := g.GetCurrentPoint()
	if err != nil {
		return nil, err
	}

	hierarchy := g.currentWorld.GetHierarchy()
	connections := make([]*entities.ConnectionInfo, 0, len(currentPoint.Connections))

	// Находим текущий субкластер
	currentSubCluster := hierarchy.FindSubClusterByPoint(currentPoint.ID)

	for _, connectionID := range currentPoint.Connections {
		point := hierarchy.FindPoint(connectionID)
		if point != nil {
			// Создаем мапу имен соединений
			connectionNames := make(map[string]string)
			for _, connID := range point.Connections {
				if connPoint := hierarchy.FindPoint(connID); connPoint != nil {
					connectionNames[connID] = connPoint.Name
				}
			}

			connectionInfo := &entities.ConnectionInfo{
				ID:              point.ID,
				Name:            point.Name,
				Description:     point.Description,
				SubClusterID:    point.SubClusterID,
				Type:            point.Type,
				Connections:     point.Connections,
				ConnectionNames: connectionNames,
				NPCs:            point.NPCs,
				IsEntryPoint:    point.IsEntryPoint,
			}

			// Определяем, является ли это переходом между субкластерами
			if currentSubCluster != nil && point.SubClusterID != currentSubCluster.ID {
				connectionInfo.IsInterCluster = true

				// Находим целевой субкластер
				targetSubCluster := hierarchy.FindSubCluster(point.SubClusterID)
				if targetSubCluster != nil {
					connectionInfo.TargetSubCluster = targetSubCluster.Name
					connectionInfo.DisplayName = fmt.Sprintf("Перейти в %s", targetSubCluster.Name)
				} else {
					connectionInfo.DisplayName = point.Name
				}
			} else {
				connectionInfo.IsInterCluster = false
				connectionInfo.DisplayName = point.Name
			}

			connections = append(connections, connectionInfo)
		}
	}

	return connections, nil
}

func (g *GameEngine) GetNPCsForCurrentPoint() ([]*entities.NPC, error) {
	currentPoint, err := g.GetCurrentPoint()
	if err != nil {
		return nil, err
	}

	npcs := make([]*entities.NPC, 0, len(currentPoint.NPCs))
	for _, npcID := range currentPoint.NPCs {
		if npc, exists := g.currentWorld.NPCs[npcID]; exists {
			npcs = append(npcs, npc)
		}
	}

	return npcs, nil
}

func (g *GameEngine) GetInteractionsForCurrentPoint() ([]entities.Interaction, error) {
	if g.gameState == nil {
		return nil, fmt.Errorf("game state not initialized")
	}

	currentPointID := g.gameState.GetCurrentPointID()
	if currentPointID == "" {
		return []entities.Interaction{}, nil
	}

	pointState := g.gameState.GetPointState(currentPointID)
	return pointState.Interactions, nil
}
