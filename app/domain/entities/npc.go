package entities

// NPC represents a non-player character
type NPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Race        string `json:"race"`
	LocationID  string `json:"location_id"`
	Description string `json:"description"`
}

// ToFrontendInfo converts NPC to frontend-compatible format (without LocationID)
func (npc *NPC) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          npc.ID,
		"name":        npc.Name,
		"race":        npc.Race,
		"description": npc.Description,
	}
}
