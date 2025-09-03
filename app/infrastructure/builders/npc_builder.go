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
	race := nb.selectRaceForLocation(clusterType, rng)
	name := nb.generateNameForRace(race, rng)
	description := nb.generateDescriptionForRace(race, rng)
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_npc_%d", point.ID, npcIndex),
		Name:        name,
		Race:        race,
		LocationID:  point.ID,
		Description: description,
	}
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
