package entities

type NPC struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Race         string       `json:"race"`
	LocationID   string       `json:"location_id"`
	Description  string       `json:"description"`
	TemperTraits TemperTraits `json:"temper_traits"`
}

func (npc *NPC) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          npc.ID,
		"name":        npc.Name,
		"race":        npc.Race,
		"description": npc.Description,
	}
}

func NewNPCFromMap(npcData interface{}) *NPC {
	npcMap, ok := npcData.(map[string]interface{})
	if !ok {
		return nil
	}
	return &NPC{
		ID:          getStringFromMap(npcMap, "id"),
		Name:        getStringFromMap(npcMap, "name"),
		Race:        getStringFromMap(npcMap, "race"),
		LocationID:  getStringFromMap(npcMap, "location_id"),
		Description: getStringFromMap(npcMap, "description"),
	}
}
