package app

import (
	"context"
	"fmt"

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
func (a *App) GetCurrentLocation() (*game.LocationInfo, error) {
	return a.gameEngine.GetCurrentLocationInfo()
}
