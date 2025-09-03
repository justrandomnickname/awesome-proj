package builders

import (
	"awesome-proj/app/domain/entities"
	"fmt"
	"math/rand"
)

type NPCBuilder struct {
	raceNames       []string
	humanNames      []string
	descriptions    map[string][]string
	locationRaceMap map[string]string
	traitSystem     *TraitSystem
}

func NewNPCBuilder() *NPCBuilder {
	return &NPCBuilder{
		raceNames: []string{"human"},
		humanNames: []string{
			"Артур", "Гвендолин", "Роланд", "Изабелла", "Гарет", "Элеонора",
			"Торвальд", "Бригитта", "Алрик", "Катарина", "Дункан", "Морген",
		},
		descriptions: map[string][]string{
			"human": {
				"Обычный человек",
			},
		},
		locationRaceMap: map[string]string{
			"village": "human",
		},
		traitSystem: NewTraitSystem(),
	}
}

func (nb *NPCBuilder) GenerateNPCsForPoint(point *entities.Point, clusterType string, rng *rand.Rand) []*entities.NPC {
	npcCount := 1 + rng.Intn(3)
	npcs := make([]*entities.NPC, 0, npcCount)

	for i := 0; i < npcCount; i++ {
		npc := nb.generateSingleNPCForPoint(point, clusterType, rng, i+1)
		npcs = append(npcs, npc)
		point.NPCs = append(point.NPCs, npc.ID)
	}
	return npcs
}

func (nb *NPCBuilder) generateSingleNPCForPoint(point *entities.Point, clusterType string, rng *rand.Rand, npcIndex int) *entities.NPC {
	id := fmt.Sprintf("%s_npc_%d", point.ID, npcIndex)
	race := nb.selectRaceForLocation(clusterType, rng)
	name := nb.generateNameForRace(race, rng)
	description := nb.generateDescriptionForRace(race, rng)
	temperTraits := nb.generateTemperTraits(race, rng, id)

	npc := &entities.NPC{
		ID:           id,
		Name:         name,
		Race:         race,
		LocationID:   point.ID,
		Description:  description,
		TemperTraits: temperTraits,
		TraitIDs:     []string{}, // Инициализируем пустым списком
	}

	// Генерируем трейты для NPC
	traits := nb.traitSystem.GenerateTraitsForNPC(npc, rng)

	// Применяем эффекты трейтов к характеристикам
	nb.traitSystem.ApplyTraitEffects(npc, traits)

	// Сохраняем ID трейтов
	for _, trait := range traits {
		npc.TraitIDs = append(npc.TraitIDs, trait.ID)
	}

	return npc
}

func (nb *NPCBuilder) selectRaceForLocation(locationType string, rng *rand.Rand) string {
	if primaryRace, exists := nb.locationRaceMap[locationType]; exists {
		if rng.Intn(100) < 80 {
			return primaryRace
		}
	}

	return nb.raceNames[rng.Intn(len(nb.raceNames))]
}

func (nb *NPCBuilder) generateNameForRace(race string, rng *rand.Rand) string {
	switch race {
	case "human":
		return nb.humanNames[rng.Intn(len(nb.humanNames))]
	default:
		return "Неизвестный"
	}
}

func (nb *NPCBuilder) generateDescriptionForRace(race string, rng *rand.Rand) string {
	if descs, exists := nb.descriptions[race]; exists {
		return descs[rng.Intn(len(descs))]
	}

	return fmt.Sprintf("Загадочный представитель расы %s", race)
}

func (nb *NPCBuilder) generateTemperTraits(race string, rng *rand.Rand, npcId string) entities.TemperTraits {
	if race == "" {
		race = "human"
	}

	return entities.TemperTraits{
		ID:             fmt.Sprintf("%s_temper_traits", npcId),
		NPC_ID:         npcId,
		Prudence:       nb.generateTraitValue(rng),
		Emotionality:   nb.generateTraitValue(rng),
		Independence:   nb.generateTraitValue(rng),
		Optimism:       nb.generateTraitValue(rng),
		Flexibility:    nb.generateTraitValue(rng),
		Aggressiveness: nb.generateTraitValue(rng),
	}
}

// generateTraitValue генерирует значение черты характера от -10 до 10
// с нормальным распределением, где экстремальные значения очень редки
func (nb *NPCBuilder) generateTraitValue(rng *rand.Rand) int {
	// Создаем весовую систему для значений от -10 до 10
	// Веса определяют вероятность выпадения каждого значения
	weights := map[int]int{
		-10: 1,  // 1% - крайний экстремум
		-9:  1,  // 1% - почти экстремум
		-8:  2,  // 2% - очень высокое значение
		-7:  3,  // 3% - высокое значение
		-6:  5,  // 5% - заметное значение
		-5:  8,  // 8% - умеренно высокое
		-4:  12, // 12% - слегка выше среднего
		-3:  15, // 15% - чуть выше среднего
		-2:  18, // 18% - близко к среднему
		-1:  20, // 20% - почти среднее
		0:   22, // 22% - абсолютно среднее (самое частое)
		1:   20, // 20% - почти среднее
		2:   18, // 18% - близко к среднему
		3:   15, // 15% - чуть выше среднего
		4:   12, // 12% - слегка выше среднего
		5:   8,  // 8% - умеренно высокое
		6:   5,  // 5% - заметное значение
		7:   3,  // 3% - высокое значение
		8:   2,  // 2% - очень высокое значение
		9:   1,  // 1% - почти экстремум
		10:  1,  // 1% - крайний экстремум
	}

	// Создаем массив всех возможных значений с учетом весов
	var weightedValues []int
	for value, weight := range weights {
		for i := 0; i < weight; i++ {
			weightedValues = append(weightedValues, value)
		}
	}

	// Выбираем случайное значение из взвешенного массива
	return weightedValues[rng.Intn(len(weightedValues))]
}
