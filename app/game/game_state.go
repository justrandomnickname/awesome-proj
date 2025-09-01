package game

// GameState holds the current game state
type GameState struct {
	CurrentLocationID string `json:"current_location_id"`
	WorldSeed         int64  `json:"world_seed"` // Сид для генерации всего мира
	// В будущем добавим:
	// CurrentPlayerID   string
	// GameTime         int64
	// etc...
}

// NewGameState creates a new game state with specified seed
func NewGameState(worldSeed int64) *GameState {
	return &GameState{
		CurrentLocationID: "start", // начинаем со стартовой локации
		WorldSeed:         worldSeed,
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

// GetWorldSeed returns the world generation seed
func (gs *GameState) GetWorldSeed() int64 {
	return gs.WorldSeed
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
