package game

import "awesome-proj/app/domain/entities"

// GameState holds the current game state
type GameState struct {
	CurrentLocationID string                           `json:"current_location_id"`
	WorldSeed         int64                            `json:"world_seed"` // Сид для генерации всего мира
	LocationStates    map[string]*entities.LocationState `json:"location_states"` // Состояния всех посещенных локаций
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
		LocationStates:    make(map[string]*entities.LocationState),
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

// GetLocationState returns the state for a specific location
func (gs *GameState) GetLocationState(locationID string) *entities.LocationState {
	state, exists := gs.LocationStates[locationID]
	if !exists {
		// Создаем новое состояние для локации при первом посещении
		state = &entities.LocationState{
			LocationID:   locationID,
			Interactions: make([]entities.Interaction, 0),
			FirstVisit:   true,
		}
		gs.LocationStates[locationID] = state
	}
	return state
}

// AddInteractionToCurrentLocation adds an interaction to the current location
func (gs *GameState) AddInteractionToCurrentLocation(interaction entities.Interaction) {
	locationState := gs.GetLocationState(gs.CurrentLocationID)
	locationState.AddInteraction(interaction)
	// После первого взаимодействия больше не первое посещение
	if locationState.FirstVisit {
		locationState.FirstVisit = false
	}
}
