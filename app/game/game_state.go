package game

import "awesome-proj/app/domain/entities"

type GameState struct {
	CurrentPointID string                          `json:"current_point_id"`
	WorldSeed      int64                           `json:"world_seed"`
	PointStates    map[string]*entities.PointState `json:"point_states"`
}

func NewGameState(worldSeed int64) *GameState {
	return &GameState{
		CurrentPointID: "start_sub_1_point_1", // начинаем с входной точки
		WorldSeed:      worldSeed,
		PointStates:    make(map[string]*entities.PointState),
	}
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

func (gs *GameState) GetPointState(pointID string) *entities.PointState {
	// Инициализируем PointStates если nil (например, после загрузки из JSON)
	if gs.PointStates == nil {
		gs.PointStates = make(map[string]*entities.PointState)
	}

	state, exists := gs.PointStates[pointID]
	if !exists {
		state = &entities.PointState{
			PointID:      pointID,
			Interactions: make([]entities.Interaction, 0),
			FirstVisit:   true,
		}
		gs.PointStates[pointID] = state
	}

	return state
}

func (gs *GameState) AddInteractionToCurrentPoint(interaction entities.Interaction) {
	pointState := gs.GetPointState(gs.CurrentPointID)
	pointState.AddInteraction(interaction)
	if pointState.FirstVisit {
		pointState.FirstVisit = false
	}
}
