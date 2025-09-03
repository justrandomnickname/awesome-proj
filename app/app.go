package app

import (
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/domain/services/prompts"
	"awesome-proj/app/game"
	"context"
	"fmt"
)

type App struct {
	ctx                                context.Context
	gameEngine                         *game.GameEngine
	subclusterStaticDescriptionService *prompts.SubclusterStaticDescriptionService
}

func NewApp() *App {
	gameEngine := game.NewGameEngine()
	return &App{
		gameEngine:                         gameEngine,
		subclusterStaticDescriptionService: prompts.NewSubclusterStaticDescriptionService(gameEngine),
	}
}
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.gameEngine.Initialize(ctx)
}
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) SaveGame(saveName string) error {
	return a.gameEngine.SaveGame(saveName)
}
func (a *App) LoadGame(filename string) error {
	return a.gameEngine.LoadGame(filename)
}
func (a *App) GetSavesList() ([]entities.SaveInfo, error) {
	return a.gameEngine.GetSavesList()
}
func (a *App) DeleteSave(filename string) error {
	return a.gameEngine.DeleteSave(filename)
}
func (a *App) NewGame(seed int64) error {
	return a.gameEngine.NewGame(seed)
}
func (a *App) PerformPlayerAction(actionText string) error {
	return a.gameEngine.PerformPlayerAction(actionText)
}
func (a *App) GetLocationHierarchy() (*entities.LocationHierarchy, error) {
	return a.gameEngine.GetLocationHierarchy()
}
func (a *App) GetCurrentPoint() (*entities.Point, error) {
	return a.gameEngine.GetCurrentPoint()
}
func (a *App) MoveToPoint(pointID string) error {
	return a.gameEngine.MoveToPoint(pointID)
}
func (a *App) GetAvailableConnections() ([]*entities.Point, error) {
	return a.gameEngine.GetAvailableConnections()
}

func (a *App) GetAvailableConnectionsInfo() ([]*entities.ConnectionInfo, error) {
	return a.gameEngine.GetAvailableConnectionsInfo()
}

func (a *App) GetNPCsForCurrentPoint() ([]*entities.NPC, error) {
	return a.gameEngine.GetNPCsForCurrentPoint()
}

func (a *App) GetInteractionsForCurrentPoint() ([]entities.Interaction, error) {
	return a.gameEngine.GetInteractionsForCurrentPoint()
}

func (a *App) GenerateSubclusterDescriptionPrompt() (string, error) {
	return a.subclusterStaticDescriptionService.GenerateAndPrintSubclusterDescriptionPrompt(a.ctx)
}
