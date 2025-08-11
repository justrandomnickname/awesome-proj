package entities

// NPC represents a non-player character
type NPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Race        string `json:"race"`
	LocationID  string `json:"location_id"`
	Description string `json:"description"`
}
