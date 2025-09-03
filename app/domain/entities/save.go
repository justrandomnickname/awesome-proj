package entities

import "time"

type SaveData struct {
	SaveName  string      `json:"save_name"`
	CreatedAt time.Time   `json:"created_at"`
	World     interface{} `json:"world"`
	GameState interface{} `json:"game_state"`
	Version   string      `json:"version"`
}

type SaveInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Filename  string    `json:"filename"`
}

func (si *SaveInfo) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":       si.Name,
		"created_at": si.CreatedAt,
		"filename":   si.Filename,
	}
}
