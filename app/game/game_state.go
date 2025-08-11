package game

// GameState holds the current game state
type GameState struct {
	CurrentLocationID string `json:"current_location_id"`
	// В будущем добавим:
	// CurrentPlayerID   string
	// GameTime         int64
	// etc...
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	return &GameState{
		CurrentLocationID: "start", // начинаем со стартовой локации
	}
}

// GetCurrentLocationID returns the current location ID
func (gs *GameState) GetCurrentLocationID() string {
	return gs.CurrentLocationID
}

// SetCurrentLocationID sets the current location ID
func (gs *GameState) SetCurrentLocationID(locationID string) {
	gs.CurrentLocationID = locationID
}

// LocationInfo represents location information for frontend
type LocationInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NPCs        []NPCInfo `json:"npcs"`
}

// NPCInfo represents NPC information for frontend
type NPCInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Race        string `json:"race"`
	Description string `json:"description"`
}
