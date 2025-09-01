package entities

import "time"

// SaveInfo represents save file metadata
type SaveInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Filename  string    `json:"filename"`
}

// ToFrontendInfo converts SaveInfo to frontend-compatible format
func (si *SaveInfo) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":       si.Name,
		"created_at": si.CreatedAt,
		"filename":   si.Filename,
	}
}
