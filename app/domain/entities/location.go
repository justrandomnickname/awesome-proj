package entities

// Location represents a game location
type Location struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Exits       map[string]string `json:"exits"`
	NPCs        []string          `json:"npcs"` // список ID НПЦ в этой локации
}
