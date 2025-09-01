package entities

import "time"

// SaveData represents the complete game state for persistence
type SaveData struct {
	SaveName    string      `json:"save_name"`
	CreatedAt   time.Time   `json:"created_at"`
	World       interface{} `json:"world"`      // используем interface{} чтобы избежать циклических импортов
	GameState   interface{} `json:"game_state"` // используем interface{} для гибкости
	Version     string      `json:"version"`    // для будущей совместимости
}
