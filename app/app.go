package app
import (
	"context"
	"fmt"
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/game"
)
type App struct {
	ctx        context.Context
	gameEngine *game.GameEngine
}
func NewApp() *App {
	return &App{
		gameEngine: game.NewGameEngine(),
	}
}
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.gameEngine.Initialize(ctx)
}
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) GetCurrentLocation() (*entities.Location, error) {
	return a.gameEngine.GetCurrentLocationInfo()
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
