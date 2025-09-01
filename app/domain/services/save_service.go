package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"awesome-proj/app/domain/aggregates"
)

// SaveService handles game save/load operations
type SaveService struct {
	savesDir string
}

// SaveData represents the complete game state for persistence
type SaveData struct {
	SaveName    string                `json:"save_name"`
	CreatedAt   time.Time            `json:"created_at"`
	World       *aggregates.World    `json:"world"`
	GameState   interface{}          `json:"game_state"` // используем interface{} для гибкости
	Version     string               `json:"version"`    // для будущей совместимости
}

// SaveInfo represents save file metadata
type SaveInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Filename  string    `json:"filename"`
}

// NewSaveService creates a new save service
func NewSaveService() *SaveService {
	savesDir := filepath.Join(".", "saves")
	
	// Создаем папку saves если её нет
	if err := os.MkdirAll(savesDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create saves directory: %v\n", err)
	}
	
	return &SaveService{
		savesDir: savesDir,
	}
}

// SaveGame saves the current game state to a file
func (s *SaveService) SaveGame(saveName string, world *aggregates.World, gameState interface{}) error {
	if saveName == "" {
		return fmt.Errorf("save name cannot be empty")
	}
	
	// Создаем безопасное имя файла
	safeFileName := s.sanitizeFileName(saveName)
	if safeFileName == "" {
		return fmt.Errorf("invalid save name")
	}
	
	saveData := &SaveData{
		SaveName:  saveName,
		CreatedAt: time.Now(),
		World:     world,
		GameState: gameState,
		Version:   "1.0",
	}
	
	filePath := filepath.Join(s.savesDir, safeFileName+".json")
	
	// Создаем папку saves если её нет
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

// LoadGame loads game state from a file
func (s *SaveService) LoadGame(filename string) (*aggregates.World, interface{}, error) {
	filePath := filepath.Join(s.savesDir, filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open save file: %v", err)
	}
	defer file.Close()
	
	var saveData SaveData
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&saveData); err != nil {
		return nil, nil, fmt.Errorf("could not decode save data: %v", err)
	}
	
	return saveData.World, saveData.GameState, nil
}

// GetSavesList returns list of available save files
func (s *SaveService) GetSavesList() ([]SaveInfo, error) {
	// Создаем папку saves если её нет
	if err := os.MkdirAll(s.savesDir, 0755); err != nil {
		return []SaveInfo{}, fmt.Errorf("could not create saves directory: %v", err)
	}
	
	files, err := os.ReadDir(s.savesDir)
	if err != nil {
		return []SaveInfo{}, fmt.Errorf("could not read saves directory: %v", err)
	}
	
	var saves []SaveInfo
	
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			// Пытаемся прочитать метаданные файла
			saveInfo, err := s.getSaveInfo(file.Name())
			if err != nil {
				// Если не можем прочитать - создаем базовую информацию
				fileInfo, _ := file.Info()
				saveInfo = SaveInfo{
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

// DeleteSave deletes a save file
func (s *SaveService) DeleteSave(filename string) error {
	filePath := filepath.Join(s.savesDir, filename)
	
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete save file: %v", err)
	}
	
	return nil
}

// LoadFirstAvailableSave loads the first available save file
func (s *SaveService) LoadFirstAvailableSave() (*aggregates.World, interface{}, error) {
	saves, err := s.GetSavesList()
	if err != nil {
		return nil, nil, err
	}
	
	if len(saves) == 0 {
		return nil, nil, fmt.Errorf("no save files found")
	}
	
	// Загружаем первый найденный сейв
	return s.LoadGame(saves[0].Filename)
}

// getSaveInfo reads save metadata from file
func (s *SaveService) getSaveInfo(filename string) (SaveInfo, error) {
	filePath := filepath.Join(s.savesDir, filename)
	
	file, err := os.Open(filePath)
	if err != nil {
		return SaveInfo{}, err
	}
	defer file.Close()
	
	var saveData SaveData
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&saveData); err != nil {
		return SaveInfo{}, err
	}
	
	return SaveInfo{
		Name:      saveData.SaveName,
		CreatedAt: saveData.CreatedAt,
		Filename:  filename,
	}, nil
}

// sanitizeFileName creates a safe filename from save name
func (s *SaveService) sanitizeFileName(saveName string) string {
	// Убираем опасные символы и заменяем пробелы
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
	
	// Добавляем timestamp для уникальности
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s", safe, timestamp)
}
