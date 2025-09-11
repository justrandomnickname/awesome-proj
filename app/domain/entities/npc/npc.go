package npc

type NPC struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Race         string       `json:"race"`
	LocationID   string       `json:"location_id"`
	Description  string       `json:"description"`
	TemperTraits TemperTraits `json:"temper_traits"`
	TraitIDs     []string     `json:"trait_ids"` // Список ID трейтов
}

func (npc *NPC) ToFrontendInfo() any {
	return map[string]any{
		"id":          npc.ID,
		"name":        npc.Name,
		"race":        npc.Race,
		"description": npc.Description,
	}
}

// GetNPCTraits возвращает список названий трейтов для использования в промптах
// Требует TraitSystem для получения названий по ID
func (npc *NPC) GetNPCTraits(traitSystem any) []string {
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

	npcEntity := &NPC{
		ID:          getStringFromMap(npcMap, "id"),
		Name:        getStringFromMap(npcMap, "name"),
		Race:        getStringFromMap(npcMap, "race"),
		LocationID:  getStringFromMap(npcMap, "location_id"),
		Description: getStringFromMap(npcMap, "description"),
	}

	// Загружаем TemperTraits
	if temperTraitsData, ok := npcMap["temper_traits"].(map[string]interface{}); ok {
		npcEntity.TemperTraits = TemperTraits{
			ID:             getStringFromMap(temperTraitsData, "id"),
			NPC_ID:         getStringFromMap(temperTraitsData, "npc_id"),
			Prudence:       getIntFromMap(temperTraitsData, "prudence"),
			Emotionality:   getIntFromMap(temperTraitsData, "emotionality"),
			Independence:   getIntFromMap(temperTraitsData, "independence"),
			Optimism:       getIntFromMap(temperTraitsData, "optimism"),
			Flexibility:    getIntFromMap(temperTraitsData, "flexibility"),
			Aggressiveness: getIntFromMap(temperTraitsData, "aggressiveness"),
		}
	}

	// Загружаем TraitIDs
	if traitIDsData, ok := npcMap["trait_ids"].([]interface{}); ok {
		npcEntity.TraitIDs = make([]string, 0, len(traitIDsData))
		for _, traitID := range traitIDsData {
			if str, ok := traitID.(string); ok {
				npcEntity.TraitIDs = append(npcEntity.TraitIDs, str)
			}
		}
	}

	return npcEntity
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if floatVal, ok := val.(float64); ok {
			return int(floatVal)
		}
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0
}
