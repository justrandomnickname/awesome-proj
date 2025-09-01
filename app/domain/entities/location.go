package entities

// Location represents a game location
type Location struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	CurrentState  string            `json:"current_state"`
	Type          string            `json:"type"`
	Exits         map[string]string `json:"exits"`
	NPCs          []string          `json:"npcs"`
	
	// Поля для фронтенда (заполняются через ToFrontendFormat)
	NPCsDetailed  []*NPC        `json:"npcs_detailed,omitempty"`
	Interactions  []Interaction `json:"interactions,omitempty"`
}

// LocationState represents the current state of a location based on interactions
type LocationState struct {
	LocationID   string        `json:"location_id"`
	Interactions []Interaction `json:"interactions"`
	// Флаг, указывающий, была ли локация уже посещена
	FirstVisit   bool          `json:"first_visit"`
}

// ToFrontendFormat заполняет поля для фронтенда
func (l *Location) ToFrontendFormat(npcs []*NPC, interactions []Interaction) {
	l.NPCsDetailed = npcs
	l.Interactions = interactions
}

// GetPlayerActions returns only player actions for this location
func (ls *LocationState) GetPlayerActions() []Interaction {
	var actions []Interaction
	for _, interaction := range ls.Interactions {
		if interaction.Type == InteractionTypePlayerAction {
			actions = append(actions, interaction)
		}
	}
	return actions
}

// GetLocationResponses returns only location responses for this location
func (ls *LocationState) GetLocationResponses() []Interaction {
	var responses []Interaction
	for _, interaction := range ls.Interactions {
		if interaction.Type == InteractionTypeLocationResponse {
			responses = append(responses, interaction)
		}
	}
	return responses
}

// AddInteraction adds a new interaction to the location state
func (ls *LocationState) AddInteraction(interaction Interaction) {
	ls.Interactions = append(ls.Interactions, interaction)
}
