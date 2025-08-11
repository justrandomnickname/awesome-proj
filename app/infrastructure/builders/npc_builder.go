package builders

import (
	"fmt"
	"math/rand"

	"awesome-proj/app/domain/entities"
)

// NPCBuilder handles NPC generation for locations
type NPCBuilder struct {
	raceNames        []string
	humanNames       []string
	dwarfNames       []string
	skavenNames      []string
	descriptions     map[string][]string
	locationRaceMap  map[string]string // какая раса основная для типа локации
}

// NewNPCBuilder creates a new NPC builder
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
			"mountain": "dwarf",  // 80% дварфы в горах
			"ruins":    "skaven", // 80% скавены в руинах  
			"village":  "human",  // 80% люди в деревнях
		},
	}
}

// GenerateNPCsForLocation generates random NPCs for a location
func (nb *NPCBuilder) GenerateNPCsForLocation(location *entities.Location, rng *rand.Rand) []*entities.NPC {
	// Generate 3-5 NPCs per location
	npcCount := 3 + rng.Intn(3) // 3-5 NPCs
	
	npcs := make([]*entities.NPC, 0, npcCount)
	
	for i := 0; i < npcCount; i++ {
		npc := nb.generateSingleNPC(location, rng, i+1)
		npcs = append(npcs, npc)
		
		// Add NPC ID to location
		location.NPCs = append(location.NPCs, npc.ID)
	}
	
	return npcs
}

// generateSingleNPC creates a single NPC for a location
func (nb *NPCBuilder) generateSingleNPC(location *entities.Location, rng *rand.Rand, npcIndex int) *entities.NPC {
	// Determine race based on location with 80% chance for primary race
	race := nb.selectRaceForLocation(location.Type, rng)
	
	// Generate name based on race
	name := nb.generateNameForRace(race, rng)
	
	// Generate description
	description := nb.generateDescriptionForRace(race, rng)
	
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_npc_%d", location.ID, npcIndex),
		Name:        name,
		Race:        race,
		LocationID:  location.ID,
		Description: description,
	}
}

// selectRaceForLocation determines NPC race based on location type
func (nb *NPCBuilder) selectRaceForLocation(locationType string, rng *rand.Rand) string {
	// Check if location has primary race
	if primaryRace, exists := nb.locationRaceMap[locationType]; exists {
		// 80% chance for primary race
		if rng.Intn(100) < 80 {
			return primaryRace
		}
	}
	
	// 20% chance or fallback - random race
	return nb.raceNames[rng.Intn(len(nb.raceNames))]
}

// generateNameForRace generates a name based on race
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

// generateDescriptionForRace generates description based on race
func (nb *NPCBuilder) generateDescriptionForRace(race string, rng *rand.Rand) string {
	if descs, exists := nb.descriptions[race]; exists {
		return descs[rng.Intn(len(descs))]
	}
	return fmt.Sprintf("Загадочный представитель расы %s", race)
}
