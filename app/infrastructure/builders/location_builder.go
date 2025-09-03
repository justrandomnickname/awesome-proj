package builders

import (
	"awesome-proj/app/domain/entities"
	"awesome-proj/app/infrastructure/builders/static"
	"fmt"
	"math/rand"
)

type WorldInterface interface {
	GetLocations() map[string]*entities.Location
	AddLocation(id string, location *entities.Location)
	AddNPC(npc *entities.NPC)
}

type ClusterTemplate struct {
	Name        string
	Type        string
	Description string
	SubCluster  SubClusterTemplate
}

type SubClusterTemplate struct {
	Name        string
	Description string
	Points      []PointTemplate
}

type PointTemplate struct {
	Name         string
	Description  string
	Type         entities.PointType
	IsEntryPoint bool
}

type LocationBuilder struct {
	locationTypes    []string
	locationNames    []string
	descriptions     map[string][]string
	npcBuilder       *NPCBuilder
	tavernBuilder    *static.TavernBuilder
	clusterTemplates []ClusterTemplate
}

func NewLocationBuilder() *LocationBuilder {
	npcBuilder := NewNPCBuilder()
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
		npcBuilder:    npcBuilder,
		tavernBuilder: static.NewTavernBuilder(),
		clusterTemplates: []ClusterTemplate{
			{
				Name: "Темный лес", Type: "forest",
				Description: "Мрачный лес, полный тайн и опасностей",
				SubCluster: SubClusterTemplate{
					Name: "Сердце леса", Description: "Самая глубокая часть темного леса",
					Points: []PointTemplate{
						{Name: "Лесная поляна", Description: "Открытая поляна среди темных деревьев", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Древний дуб", Description: "Огромный дуб, помнящий века", Type: entities.PointTypeSpecial},
						{Name: "Заросшая тропа", Description: "Едва заметная тропа, ведущая вглубь", Type: entities.PointTypeRegular},
					},
				},
			},
			{
				Name: "Забытые руины", Type: "ruins",
				Description: "Остатки древней цивилизации",
				SubCluster: SubClusterTemplate{
					Name: "Развалины башни", Description: "Полуразрушенная башня неизвестного назначения",
					Points: []PointTemplate{
						{Name: "Входная арка", Description: "Разрушенный вход в башню", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Винтовая лестница", Description: "Обвалившаяся лестница, ведущая наверх", Type: entities.PointTypeRegular},
						{Name: "Комната с фресками", Description: "Зал с потускневшими древними росписями", Type: entities.PointTypeSpecial},
					},
				},
			},
			{
				Name: "Старая пещера", Type: "cave",
				Description: "Глубокая пещера с эхом прошлого",
				SubCluster: SubClusterTemplate{
					Name: "Подземные гроты", Description: "Система соединенных пещер",
					Points: []PointTemplate{
						{Name: "Вход в пещеру", Description: "Широкий вход, откуда веет прохладой", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Сталактитовый зал", Description: "Большая полость с острыми сталактитами", Type: entities.PointTypeRegular},
						{Name: "Подземное озеро", Description: "Темное озеро в глубине пещеры", Type: entities.PointTypeSpecial},
					},
				},
			},
			{
				Name: "Заброшенная деревня", Type: "village",
				Description: "Пустая деревня, покинутая жителями",
				SubCluster: SubClusterTemplate{
					Name: "Центр деревни", Description: "Главная площадь заброшенной деревни",
					Points: []PointTemplate{
						{Name: "Деревенская площадь", Description: "Пустая площадь с сухим колодцем", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Старый дом", Description: "Полуразрушенный дом с провалившейся крышей", Type: entities.PointTypeRegular},
						{Name: "Церковь", Description: "Маленькая церковь с разбитыми окнами", Type: entities.PointTypeSpecial},
					},
				},
			},
			{
				Name: "Мрачные болота", Type: "swamp",
				Description: "Топкие болота с ядовитыми испарениями",
				SubCluster: SubClusterTemplate{
					Name: "Трясина", Description: "Самая опасная часть болот",
					Points: []PointTemplate{
						{Name: "Гать", Description: "Деревянная дорога через болото", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Болотный островок", Description: "Небольшой сухой островок среди трясины", Type: entities.PointTypeRegular},
						{Name: "Хижина на сваях", Description: "Заброшенная хижина болотника", Type: entities.PointTypeSpecial},
					},
				},
			},
			{
				Name: "Скалистый утес", Type: "mountain",
				Description: "Высокий утес с видом на окрестности",
				SubCluster: SubClusterTemplate{
					Name: "Горная тропа", Description: "Опасная тропа по краю утеса",
					Points: []PointTemplate{
						{Name: "Подножие утеса", Description: "Начало подъема на скалистый утес", Type: entities.PointTypeEntry, IsEntryPoint: true},
						{Name: "Горный карниз", Description: "Узкий выступ в скале", Type: entities.PointTypeRegular},
						{Name: "Вершина утеса", Description: "Высшая точка с панорамным видом", Type: entities.PointTypeSpecial},
					},
				},
			},
		},
	}
}

func (lb *LocationBuilder) GenerateLocationHierarchy(rng *rand.Rand) *entities.LocationHierarchy {
	hierarchy := entities.NewLocationHierarchy()

	template := lb.clusterTemplates[rng.Intn(len(lb.clusterTemplates))]
	cluster := lb.buildClusterFromTemplate(template, "start", rng)
	hierarchy.AddCluster(cluster)

	return hierarchy
}

func (lb *LocationBuilder) GenerateLocationHierarchyWithNPCs(world WorldInterface, rng *rand.Rand) *entities.LocationHierarchy {
	hierarchy := entities.NewLocationHierarchy()

	// Выбираем рандомный шаблон кластера
	template := lb.clusterTemplates[rng.Intn(len(lb.clusterTemplates))]
	cluster := lb.buildHybridClusterWithTavern(template, "start", rng, world)
	hierarchy.AddCluster(cluster)

	return hierarchy
}

func (lb *LocationBuilder) buildClusterFromTemplate(template ClusterTemplate, clusterID string, rng *rand.Rand) *entities.Cluster {
	cluster := &entities.Cluster{
		ID:          clusterID,
		Name:        template.Name,
		Description: template.Description,
		Type:        template.Type,
		SubClusters: make(map[string]*entities.SubCluster),
	}

	subClusterID := fmt.Sprintf("%s_sub_1", clusterID)
	subCluster := lb.buildSubClusterFromTemplate(template.SubCluster, subClusterID, clusterID, rng)
	cluster.AddSubCluster(subCluster)

	if len(subCluster.EntryPoints) > 0 {
		cluster.MainPoint = subCluster.EntryPoints[0]
	}

	return cluster
}

func (lb *LocationBuilder) buildClusterFromTemplateWithNPCs(template ClusterTemplate, clusterID string, rng *rand.Rand, world WorldInterface) *entities.Cluster {
	cluster := &entities.Cluster{
		ID:          clusterID,
		Name:        template.Name,
		Description: template.Description,
		Type:        template.Type,
		SubClusters: make(map[string]*entities.SubCluster),
	}

	subClusterID := fmt.Sprintf("%s_sub_1", clusterID)
	subCluster := lb.buildSubClusterFromTemplateWithNPCs(template.SubCluster, subClusterID, clusterID, template.Type, rng, world)
	cluster.AddSubCluster(subCluster)

	if len(subCluster.EntryPoints) > 0 {
		cluster.MainPoint = subCluster.EntryPoints[0]
	}

	return cluster
}

func (lb *LocationBuilder) buildHybridClusterWithTavern(template ClusterTemplate, clusterID string, rng *rand.Rand, world WorldInterface) *entities.Cluster {
	cluster := &entities.Cluster{
		ID:          clusterID,
		Name:        template.Name,
		Description: template.Description,
		Type:        template.Type,
		SubClusters: make(map[string]*entities.SubCluster),
	}

	// Создаем рандомный субкластер из шаблона
	randomSubClusterID := fmt.Sprintf("%s_sub_random", clusterID)
	randomSubCluster := lb.buildSubClusterFromTemplateWithNPCs(template.SubCluster, randomSubClusterID, clusterID, template.Type, rng, world)
	cluster.AddSubCluster(randomSubCluster)

	// Получаем entry point рандомного субкластера для подключения к таверне
	var exitPointID string
	if len(randomSubCluster.EntryPoints) > 0 {
		exitPointID = randomSubCluster.EntryPoints[0]
	}

	// Создаем статичный субкластер "Таверна" с подключением к рандомному субкластеру
	tavernSubCluster := lb.tavernBuilder.GenerateTavernDefault(clusterID, exitPointID, world)
	cluster.AddSubCluster(tavernSubCluster)

	// Добавляем обратное соединение от рандомного субкластера к таверне
	if len(tavernSubCluster.EntryPoints) > 0 && exitPointID != "" {
		tavernEntryPointID := tavernSubCluster.EntryPoints[0]
		if randomEntryPoint, exists := randomSubCluster.Points[exitPointID]; exists {
			randomEntryPoint.AddConnection(tavernEntryPointID)
		}
	}

	// Устанавливаем главный вход через таверну
	if len(tavernSubCluster.EntryPoints) > 0 {
		cluster.MainPoint = tavernSubCluster.EntryPoints[0]
	}

	return cluster
}

func (lb *LocationBuilder) buildSubClusterFromTemplate(template SubClusterTemplate, subClusterID, clusterID string, rng *rand.Rand) *entities.SubCluster {
	subCluster := &entities.SubCluster{
		ID:          subClusterID,
		Name:        template.Name,
		Description: template.Description,
		ClusterID:   clusterID,
		EntryPoints: make([]string, 0),
		Points:      make(map[string]*entities.Point),
	}

	for i, pointTemplate := range template.Points {
		pointID := fmt.Sprintf("%s_point_%d", subClusterID, i+1)
		point := &entities.Point{
			ID:           pointID,
			Name:         pointTemplate.Name,
			Description:  pointTemplate.Description,
			SubClusterID: subClusterID,
			Type:         pointTemplate.Type,
			Connections:  make([]string, 0),
			NPCs:         make([]string, 0),
			IsEntryPoint: pointTemplate.IsEntryPoint,
		}
		subCluster.AddPoint(point)
	}

	lb.connectPointsInSubCluster(subCluster)

	return subCluster
}

func (lb *LocationBuilder) buildSubClusterFromTemplateWithNPCs(template SubClusterTemplate, subClusterID, clusterID, clusterType string, rng *rand.Rand, world WorldInterface) *entities.SubCluster {
	subCluster := &entities.SubCluster{
		ID:          subClusterID,
		Name:        template.Name,
		Description: template.Description,
		ClusterID:   clusterID,
		EntryPoints: make([]string, 0),
		Points:      make(map[string]*entities.Point),
	}

	for i, pointTemplate := range template.Points {
		pointID := fmt.Sprintf("%s_point_%d", subClusterID, i+1)
		point := &entities.Point{
			ID:           pointID,
			Name:         pointTemplate.Name,
			Description:  pointTemplate.Description,
			SubClusterID: subClusterID,
			Type:         pointTemplate.Type,
			Connections:  make([]string, 0),
			NPCs:         make([]string, 0),
			IsEntryPoint: pointTemplate.IsEntryPoint,
		}
		subCluster.AddPoint(point)

		// Генерируем NPCs для каждого Point
		npcs := lb.npcBuilder.GenerateNPCsForPoint(point, clusterType, rng)
		for _, npc := range npcs {
			world.AddNPC(npc)
		}
	}

	lb.connectPointsInSubCluster(subCluster)

	return subCluster
}

func (lb *LocationBuilder) connectPointsInSubCluster(subCluster *entities.SubCluster) {
	points := make([]*entities.Point, 0, len(subCluster.Points))
	for _, point := range subCluster.Points {
		points = append(points, point)
	}

	for i, point := range points {
		for j, otherPoint := range points {
			if i != j {
				point.AddConnection(otherPoint.ID)
			}
		}
	}
}
