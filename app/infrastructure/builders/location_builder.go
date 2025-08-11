package builders

import (
	"fmt"
	"math/rand"

	"awesome-proj/app/domain/entities"
)

// WorldInterface defines what builder needs from World
type WorldInterface interface {
	GetLocations() map[string]*entities.Location
	AddLocation(id string, location *entities.Location)
	AddNPC(npc *entities.NPC)
}

// LocationBuilder handles location generation for worlds
type LocationBuilder struct {
	locationTypes []string
	locationNames []string
	descriptions  map[string][]string
	npcBuilder    *NPCBuilder
}

// NewLocationBuilder creates a new location builder
func NewLocationBuilder() *LocationBuilder {
	return &LocationBuilder{
		locationTypes: []string{"forest", "cave", "village", "ruins", "swamp", "mountain"},
		locationNames: []string{
			"Темный лес", "Забытые руины", "Старая пещера",
			"Заброшенная деревня", "Мрачные болота", "Скалистый утес",
			"Древний храм", "Глухая чаща", "Каменные врата",
			"Проклятый алтарь", "Железный рудник", "Кладбище волков",
		},
		descriptions: map[string][]string{
			"forest": {
				"Густой темный лес, где ветви переплетаются над головой",
				"Древний лес, полный тайн и опасностей",
				"Мрачная чаща с едва заметными тропами",
				"Заколдованный лес, где шепчутся деревья",
			},
			"cave": {
				"Сырая пещера с эхом капающей воды",
				"Глубокая пещера с множеством ходов",
				"Темная пещера, откуда веет холодом",
				"Пещера с блестящими кристаллами на стенах",
			},
			"village": {
				"Заброшенная деревня с полуразрушенными домами",
				"Пустая деревня, где никто не живет уже много лет",
				"Старая деревня с призрачной атмосферой",
				"Деревня, покинутая жителями в спешке",
			},
			"ruins": {
				"Древние руины, поросшие плющом",
				"Остатки величественного замка",
				"Руины храма неизвестного божества",
				"Разрушенная башня мага",
			},
			"swamp": {
				"Мрачное болото с ядовитыми испарениями",
				"Топкие болота, где легко заблудиться",
				"Болото, где обитают странные создания",
				"Затопленный лес с торчащими стволами",
			},
			"mountain": {
				"Скалистая вершина, овеваемая ветрами",
				"Горный перевал между острыми пиками",
				"Пещера в горном склоне",
				"Плато на вершине горы",
			},
		},
		npcBuilder: NewNPCBuilder(),
	}
}

// GenerateRandomLocations creates random locations for a world
func (lb *LocationBuilder) GenerateRandomLocations(world WorldInterface, rng *rand.Rand, count int) {
	// Create one starting location as a normal generated location
	usedNames := make(map[string]bool)
	startLocation := lb.buildRandomLocation(rng, "start", usedNames)
	world.AddLocation("start", startLocation)
	
	// Generate NPCs for starting location
	npcs := lb.npcBuilder.GenerateNPCsForLocation(startLocation, rng)
	for _, npc := range npcs {
		world.AddNPC(npc)
	}
}

// buildRandomLocation creates a single random location
func (lb *LocationBuilder) buildRandomLocation(rng *rand.Rand, locationID string, usedNames map[string]bool) *entities.Location {
	// Pick random type and name
	locType := lb.locationTypes[rng.Intn(len(lb.locationTypes))]
	var name string
	for {
		name = lb.locationNames[rng.Intn(len(lb.locationNames))]
		if !usedNames[name] {
			usedNames[name] = true
			break
		}
	}
	
	// Pick description for this type
	var description string
	if descs, exists := lb.descriptions[locType]; exists {
		description = descs[rng.Intn(len(descs))]
	} else {
		description = fmt.Sprintf("Загадочное место типа %s", locType)
	}
	
	return &entities.Location{
		ID:          locationID,
		Name:        name,
		Description: description,
		Type:        locType,
		Exits:       make(map[string]string),
		NPCs:        make([]string, 0),
	}
}
