package entities

import "time"

type InteractionType string

const (
	InteractionTypePlayerAction     InteractionType = "player_action"
	InteractionTypeLocationResponse InteractionType = "location_response"
	InteractionTypeLocationState    InteractionType = "location_state"
	InteractionTypePlayerMovement   InteractionType = "player_movement"
)

type Interaction struct {
	ID                string          `json:"id"`
	Type              InteractionType `json:"type"`
	Content           string          `json:"content"`
	AdditionalContent string          `json:"additional_content,omitempty"`
	LocationID        string          `json:"location_id"`
	Timestamp         time.Time       `json:"timestamp"`
}

func (i *Interaction) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":                 i.ID,
		"type":               string(i.Type),
		"content":            i.Content,
		"timestamp":          i.Timestamp.Format("02.01 15:04:05"),
		"additional_content": i.AdditionalContent,
	}
}
