package app

import (
	"context"
	"fmt"

	"awesome-proj/app/domain/entities"
	"awesome-proj/app/game"
)

// App struct - основная структура приложения
type App struct {
	ctx        context.Context
	gameEngine *game.GameEngine
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		gameEngine: game.NewGameEngine(),
	}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	// Initialize game engine
	a.gameEngine.Initialize(ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetCurrentLocation returns current location info for frontend
func (a *App) GetCurrentLocation() (*entities.Location, error) {
	return a.gameEngine.GetCurrentLocationInfo()
}

// SaveGame saves the current game with given name
func (a *App) SaveGame(saveName string) error {
	return a.gameEngine.SaveGame(saveName)
}

// LoadGame loads a game from save file
func (a *App) LoadGame(filename string) error {
	return a.gameEngine.LoadGame(filename)
}

// GetSavesList returns list of available saves
func (a *App) GetSavesList() ([]entities.SaveInfo, error) {
	return a.gameEngine.GetSavesList()
}

// DeleteSave deletes a save file
func (a *App) DeleteSave(filename string) error {
	return a.gameEngine.DeleteSave(filename)
}

// NewGame starts a new game with specified seed
func (a *App) NewGame(seed int64) error {
	return a.gameEngine.NewGame(seed)
}

// PerformPlayerAction handles player action in current location
func (a *App) PerformPlayerAction(actionText string) error {
	return a.gameEngine.PerformPlayerAction(actionText)
}
