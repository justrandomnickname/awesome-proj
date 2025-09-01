package entities

import "time"

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
