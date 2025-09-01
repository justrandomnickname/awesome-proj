package app

import (
	"context"
	"fmt"
	"time"

	"awesome-proj/app/game"
)

// SaveInfo represents save file metadata for frontend
type SaveInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Filename  string    `json:"filename"`
}

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

// SaveGame saves the current game with given name
func (a *App) SaveGame(saveName string) error {
	return a.gameEngine.SaveGame(saveName)
}

// LoadGame loads a game from save file
func (a *App) LoadGame(filename string) error {
	return a.gameEngine.LoadGame(filename)
}

// GetSavesList returns list of available saves
func (a *App) GetSavesList() ([]SaveInfo, error) {
	engineSaves, err := a.gameEngine.GetSavesList()
	if err != nil {
		return nil, err
	}
	
	// Конвертируем из services.SaveInfo в app.SaveInfo
	// Инициализируем как пустой слайс, чтобы JSON возвращал [] вместо null
	appSaves := make([]SaveInfo, 0, len(engineSaves))
	for _, save := range engineSaves {
		appSaves = append(appSaves, SaveInfo{
			Name:      save.Name,
			CreatedAt: save.CreatedAt,
			Filename:  save.Filename,
		})
	}
	
	return appSaves, nil
}

// DeleteSave deletes a save file
func (a *App) DeleteSave(filename string) error {
	return a.gameEngine.DeleteSave(filename)
}

// NewGame starts a new game
func (a *App) NewGame() error {
	return a.gameEngine.NewGame()
}
