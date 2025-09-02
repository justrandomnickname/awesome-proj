package game

import "awesome-proj/app/domain/entities"

type GameState struct {
	CurrentLocationID string                             `json:"current_location_id"`
	CurrentPointID    string                             `json:"current_point_id"`
	WorldSeed         int64                              `json:"world_seed"`
	LocationStates    map[string]*entities.LocationState `json:"location_states"`
}

func NewGameState(worldSeed int64) *GameState {
	return &GameState{
		CurrentLocationID: "start",
		CurrentPointID:    "start_sub_1_point_1", // начинаем с входной точки
		WorldSeed:         worldSeed,
		LocationStates:    make(map[string]*entities.LocationState),
	}
}

func (gs *GameState) GetCurrentLocationID() string {
	return gs.CurrentLocationID
}

func (gs *GameState) SetCurrentLocationID(locationID string) {
	gs.CurrentLocationID = locationID
}
func (gs *GameState) GetCurrentPointID() string {
	return gs.CurrentPointID
}
func (gs *GameState) SetCurrentPointID(pointID string) {
	gs.CurrentPointID = pointID
}

func (gs *GameState) GetWorldSeed() int64 {
	return gs.WorldSeed
}

func (gs *GameState) GetLocationState(locationID string) *entities.LocationState {
	state, exists := gs.LocationStates[locationID]
	if !exists {
		state = &entities.LocationState{
			LocationID:   locationID,
			Interactions: make([]entities.Interaction, 0),
			FirstVisit:   true,
		}
		gs.LocationStates[locationID] = state
	}

	return state
}

func (gs *GameState) AddInteractionToCurrentLocation(interaction entities.Interaction) {
	locationState := gs.GetLocationState(gs.CurrentLocationID)
	locationState.AddInteraction(interaction)
	if locationState.FirstVisit {
		locationState.FirstVisit = false
	}

}
