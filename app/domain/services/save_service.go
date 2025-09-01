package services
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/domain/entities"
)
type SaveService struct {
	savesDir string
}
func NewSaveService() *SaveService {
	savesDir := filepath.Join(".", "saves")
	if err := os.MkdirAll(savesDir, 0755); err != nil {
	}
	return &SaveService{
		savesDir: savesDir,
	}
}
func (s *SaveService) SaveGame(saveName string, world *aggregates.World, gameState interface{}) error {
	if saveName == "" {
		return fmt.Errorf("save name cannot be empty")
	}
	safeFileName := s.sanitizeFileName(saveName)
	if safeFileName == "" {
		return fmt.Errorf("invalid save name")
	}
	saveData := &entities.SaveData{
		SaveName:  saveName,
		CreatedAt: time.Now(),
		World:     world,
		GameState: gameState,
		Version:   "1.0",
	}
	filePath := filepath.Join(s.savesDir, safeFileName+".json")
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return fmt.Errorf("could not create saves directory: %v", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create save file: %v", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Красивое форматирование
	if err := encoder.Encode(saveData); err != nil {
		return fmt.Errorf("could not encode save data: %v", err)
	}
	return nil
}
func (s *SaveService) LoadGame(filename string) (*aggregates.World, interface{}, error) {
	filePath := filepath.Join(s.savesDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open save file: %v", err)
	}
	defer file.Close()
	var saveData entities.SaveData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&saveData); err != nil {
		return nil, nil, fmt.Errorf("could not decode save data: %v", err)
	}
	world, err := deserializeWorldFromMap(saveData.World)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid world data in save file: %v", err)
	}
	return world, saveData.GameState, nil
}
func deserializeWorldFromMap(worldData interface{}) (*aggregates.World, error) {
	worldMap, ok := worldData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("world data is not a map")
	}
	name, _ := worldMap["name"].(string)
	var seed int64
	if seedFloat, ok := worldMap["seed"].(float64); ok {
		seed = int64(seedFloat)
	}
	world := aggregates.NewEmptyWorld(name, seed)
	if locationsData, ok := worldMap["locations"].(map[string]interface{}); ok {
		for locationID, locationData := range locationsData {
			if location := entities.NewLocationFromMap(locationData); location != nil {
				world.AddLocation(locationID, location)
			}
		}
	}
	if npcsData, ok := worldMap["npcs"].(map[string]interface{}); ok {
		for _, npcData := range npcsData {
			if npc := entities.NewNPCFromMap(npcData); npc != nil {
				world.AddNPC(npc)
			}
		}
	}
	return world, nil
}
func (s *SaveService) GetSavesList() ([]entities.SaveInfo, error) {
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return []entities.SaveInfo{}, fmt.Errorf("could not create saves directory: %v", err)
	}
	files, err := os.ReadDir(s.savesDir)
	if err != nil {
		return []entities.SaveInfo{}, fmt.Errorf("could not read saves directory: %v", err)
	}
	var saves []entities.SaveInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			saveInfo, err := s.getSaveInfo(file.Name())
			if err != nil {
				fileInfo, _ := file.Info()
				saveInfo = entities.SaveInfo{
					Name:      strings.TrimSuffix(file.Name(), ".json"),
					CreatedAt: fileInfo.ModTime(),
					Filename:  file.Name(),
				}
			}
			saves = append(saves, saveInfo)
		}
	}
	return saves, nil
}
func (s *SaveService) DeleteSave(filename string) error {
	filePath := filepath.Join(s.savesDir, filename)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete save file: %v", err)
	}
	return nil
}
func (s *SaveService) LoadFirstAvailableSave() (*aggregates.World, interface{}, error) {
	saves, err := s.GetSavesList()
	if err != nil {
		return nil, nil, err
	}
	if len(saves) == 0 {
		return nil, nil, fmt.Errorf("no save files found")
	}
	return s.LoadGame(saves[0].Filename)
}
func (s *SaveService) getSaveInfo(filename string) (entities.SaveInfo, error) {
	filePath := filepath.Join(s.savesDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return entities.SaveInfo{}, err
	}
	defer file.Close()
	var saveData entities.SaveData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&saveData); err != nil {
		return entities.SaveInfo{}, err
	}
	return entities.SaveInfo{
		Name:      saveData.SaveName,
		CreatedAt: saveData.CreatedAt,
		Filename:  filename,
	}, nil
}
func (s *SaveService) sanitizeFileName(saveName string) string {
	safe := strings.ReplaceAll(saveName, " ", "_")
	safe = strings.ReplaceAll(safe, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	safe = strings.ReplaceAll(safe, "*", "_")
	safe = strings.ReplaceAll(safe, "?", "_")
	safe = strings.ReplaceAll(safe, "\"", "_")
	safe = strings.ReplaceAll(safe, "<", "_")
	safe = strings.ReplaceAll(safe, ">", "_")
	safe = strings.ReplaceAll(safe, "|", "_")
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s", safe, timestamp)
}
