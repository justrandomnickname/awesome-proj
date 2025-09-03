package entities

type Location struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	CurrentState string            `json:"current_state"`
	Type         string            `json:"type"`
	Exits        map[string]string `json:"exits"`
	NPCs         []string          `json:"npcs"`
	NPCsDetailed []*NPC            `json:"npcs_detailed,omitempty"`
	Interactions []Interaction     `json:"interactions,omitempty"`
}

type PointState struct {
	PointID      string        `json:"point_id"`
	Interactions []Interaction `json:"interactions"`
	FirstVisit   bool          `json:"first_visit"`
}

func (l *Location) ToFrontendFormat(npcs []*NPC, interactions []Interaction) {
	l.NPCsDetailed = npcs
	l.Interactions = interactions
}

func (ps *PointState) GetPlayerActions() []Interaction {
	var actions []Interaction
	for _, interaction := range ps.Interactions {
		if interaction.Type == InteractionTypePlayerAction {
			actions = append(actions, interaction)
		}
	}

	return actions
}

func (ps *PointState) GetLocationResponses() []Interaction {
	var responses []Interaction
	for _, interaction := range ps.Interactions {
		if interaction.Type == InteractionTypeLocationResponse {
			responses = append(responses, interaction)
		}
	}

	return responses
}

func (ps *PointState) AddInteraction(interaction Interaction) {
	ps.Interactions = append(ps.Interactions, interaction)
}

func NewLocationFromMap(locationData interface{}) *Location {
	locationMap, ok := locationData.(map[string]interface{})
	if !ok {
		return nil
	}

	location := &Location{
		ID:           getStringFromMap(locationMap, "id"),
		Name:         getStringFromMap(locationMap, "name"),
		Description:  getStringFromMap(locationMap, "description"),
		CurrentState: getStringFromMap(locationMap, "current_state"),
		Type:         getStringFromMap(locationMap, "type"),
		Exits:        make(map[string]string),
		NPCs:         make([]string, 0),
	}

	if exitsData, ok := locationMap["exits"].(map[string]interface{}); ok {
		for key, value := range exitsData {
			if strValue, ok := value.(string); ok {
				location.Exits[key] = strValue
			}
		}
	}

	if npcsData, ok := locationMap["npcs"].([]interface{}); ok {
		for _, npcID := range npcsData {
			if strNpcID, ok := npcID.(string); ok {
				location.NPCs = append(location.NPCs, strNpcID)
			}
		}
	}

	return location
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}

	return ""
}
