package entities

// Location represents a game location
type Location struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	CurrentState string            `json:"current_state"`
	Type        string            `json:"type"`
	Exits       map[string]string `json:"exits"`
	NPCs        []string          `json:"npcs"`
}
