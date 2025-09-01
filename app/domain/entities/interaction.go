package entities

import "time"

// InteractionType defines the type of interaction
type InteractionType string

const (
	InteractionTypePlayerAction   InteractionType = "player_action"
	InteractionTypeLocationResponse InteractionType = "location_response"
	InteractionTypeLocationState   InteractionType = "location_state"
)

// Interaction represents a single interaction (action or response)
type Interaction struct {
	ID          string          `json:"id"`
	Type        InteractionType `json:"type"`
	Content     string          `json:"content"`
	LocationID  string          `json:"location_id"`
	Timestamp   time.Time       `json:"timestamp"`
	// В будущем можно добавить:
	// PlayerID    string          `json:"player_id"`
	// TargetNPCID string          `json:"target_npc_id,omitempty"`
}

// ToFrontendInfo converts Interaction to frontend-compatible format
func (i *Interaction) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        i.ID,
		"type":      string(i.Type),
		"content":   i.Content,
		"timestamp": i.Timestamp.Format("15:04:05"),
	}
}
