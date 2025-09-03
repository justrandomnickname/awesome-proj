package entities

type NPC struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Race         string       `json:"race"`
	LocationID   string       `json:"location_id"`
	Description  string       `json:"description"`
	TemperTraits TemperTraits `json:"temper_traits"`
	TraitIDs     []string     `json:"trait_ids"` // Список ID трейтов
}

func (npc *NPC) ToFrontendInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":          npc.ID,
		"name":        npc.Name,
		"race":        npc.Race,
		"description": npc.Description,
	}
}

// GetNPCTraits возвращает список названий трейтов для использования в промптах
// Требует TraitSystem для получения названий по ID
func (npc *NPC) GetNPCTraits(traitSystem interface{}) []string {
	// Интерфейс для избежания циклической зависимости
	// В реальном использовании будет передан *TraitSystem
	if len(npc.TraitIDs) == 0 {
		return []string{}
	}

	// Проверяем, что переданный объект имеет метод GetTraitNamesByIDs
	if ts, ok := traitSystem.(interface {
		GetTraitNamesByIDs([]string) []string
	}); ok {
		return ts.GetTraitNamesByIDs(npc.TraitIDs)
	}

	// Если метод недоступен, возвращаем пустой список
	return []string{}
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
