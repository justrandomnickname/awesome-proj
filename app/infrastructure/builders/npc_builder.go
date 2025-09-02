package builders

import (
	"awesome-proj/app/domain/entities"
	"fmt"
	"math/rand"
)

type NPCBuilder struct {
	raceNames       []string
	humanNames      []string
	dwarfNames      []string
	skavenNames     []string
	descriptions    map[string][]string
	locationRaceMap map[string]string
}

func NewNPCBuilder() *NPCBuilder {
	return &NPCBuilder{
		raceNames: []string{"human", "dwarf", "skaven"},
		humanNames: []string{
			"Артур", "Гвендолин", "Роланд", "Изабелла", "Гарет", "Элеонора",
			"Торвальд", "Бригитта", "Алрик", "Катарина", "Дункан", "Морген",
		},
		dwarfNames: []string{
			"Торин", "Дайн", "Балин", "Двалин", "Кили", "Фили",
			"Грои", "Нали", "Дори", "Ори", "Бифур", "Бофур",
		},
		skavenNames: []string{
			"Скритч", "Снеак", "Гнаурр", "Скиттер", "Твитч", "Читтер",
			"Скурк", "Нибблз", "Визкил", "Клок", "Снифф", "Ратти",
		},
		descriptions: map[string][]string{
			"human": {
				"Усталый путешественник с мешком за спиной",
				"Местный торговец, знающий все дороги",
				"Бывший солдат, ищущий приключений",
				"Мудрый старик с длинной бородой",
				"Молодая женщина с решительным взглядом",
			},
			"dwarf": {
				"Крепкий дварф с могучей киркой",
				"Бородатый мастер с молотом в руках",
				"Дварф-шахтер, весь в угольной пыли",
				"Воин клана с боевым топором",
				"Старый дварф с мудрыми глазами",
			},
			"skaven": {
				"Подозрительный крыс-человек с красными глазами",
				"Юркий скавен с острыми когтями",
				"Старый скавен-чародей с посохом",
				"Скавен-разведчик в потрепанном плаще",
				"Зловещий крысолюд с кривой ухмылкой",
			},
		},
		locationRaceMap: map[string]string{
			"mountain": "dwarf",
			"ruins":    "skaven",
			"village":  "human",
		},
	}
}

func (nb *NPCBuilder) GenerateNPCsForLocation(location *entities.Location, rng *rand.Rand) []*entities.NPC {
	npcCount := 3 + rng.Intn(3) // 3-5 NPCs
	npcs := make([]*entities.NPC, 0, npcCount)
	for i := 0; i < npcCount; i++ {
		npc := nb.generateSingleNPC(location, rng, i+1)
		npcs = append(npcs, npc)
		location.NPCs = append(location.NPCs, npc.ID)
	}
	return npcs
}

func (nb *NPCBuilder) GenerateNPCsForPoint(point *entities.Point, clusterType string, rng *rand.Rand) []*entities.NPC {
	// Для каждого поинта генерируем 1-3 NPCs
	npcCount := 1 + rng.Intn(3)
	npcs := make([]*entities.NPC, 0, npcCount)

	for i := 0; i < npcCount; i++ {
		npc := nb.generateSingleNPCForPoint(point, clusterType, rng, i+1)
		npcs = append(npcs, npc)
		point.NPCs = append(point.NPCs, npc.ID)
	}
	return npcs
}

func (nb *NPCBuilder) generateSingleNPC(location *entities.Location, rng *rand.Rand, npcIndex int) *entities.NPC {
	race := nb.selectRaceForLocation(location.Type, rng)
	name := nb.generateNameForRace(race, rng)
	description := nb.generateDescriptionForRace(race, rng)
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_npc_%d", location.ID, npcIndex),
		Name:        name,
		Race:        race,
		LocationID:  location.ID,
		Description: description,
	}
}

func (nb *NPCBuilder) generateSingleNPCForPoint(point *entities.Point, clusterType string, rng *rand.Rand, npcIndex int) *entities.NPC {
	race := nb.selectRaceForLocation(clusterType, rng)
	name := nb.generateNameForRace(race, rng)
	description := nb.generateDescriptionForRace(race, rng)
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_npc_%d", point.ID, npcIndex),
		Name:        name,
		Race:        race,
		LocationID:  point.ID, // Используем Point ID как LocationID
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
	case "dwarf":
		return nb.dwarfNames[rng.Intn(len(nb.dwarfNames))]
	case "skaven":
		return nb.skavenNames[rng.Intn(len(nb.skavenNames))]
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

func (nb *NPCBuilder) GenerateInnkeeper(pointID string) *entities.NPC {
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_innkeeper", pointID),
		Name:        "Боб Трактирщик",
		Race:        "human",
		LocationID:  pointID,
		Description: "Дружелюбный хозяин таверны с круглым животом и добрыми глазами",
	}
}
