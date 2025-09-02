package static

import (
	"awesome-proj/app/domain/entities"
	"fmt"
)

type NPCBuilderInterface interface {
	GenerateInnkeeper(pointID string) *entities.NPC
}

type WorldInterface interface {
	AddNPC(npc *entities.NPC)
}

type TavernBuilder struct {
	npcBuilder NPCBuilderInterface
}

func NewTavernBuilder(npcBuilder NPCBuilderInterface) *TavernBuilder {
	return &TavernBuilder{
		npcBuilder: npcBuilder,
	}
}

// GenerateTavernDefault создает стандартную таверну с подключением к указанной внешней точке
func (tb *TavernBuilder) GenerateTavernDefault(clusterID string, exitPointID string, world WorldInterface) *entities.SubCluster {
	subClusterID := fmt.Sprintf("%s_tavern", clusterID)

	subCluster := &entities.SubCluster{
		ID:          subClusterID,
		Name:        "Таверна",
		Description: "Уютная таверна - место отдыха путешественников",
		ClusterID:   clusterID,
		EntryPoints: make([]string, 0),
		Points:      make(map[string]*entities.Point),
	}

	// Создаем поинты таверны
	tavernPoints := []struct {
		name        string
		description string
		pointType   entities.PointType
		isEntry     bool
	}{
		{"Главный холл", "Просторный зал таверны с большим камином и деревянными столами", entities.PointTypeEntry, true},
		{"Подвал", "Прохладный подвал с бочками вина и припасами", entities.PointTypeRegular, false},
		{"Второй этаж", "Верхний этаж с номерами для постояльцев", entities.PointTypeRegular, false},
		{"Бедные номера", "Простые комнаты для небогатых путешественников", entities.PointTypeRegular, false},
		{"Богатые номера", "Роскошные апартаменты для состоятельных гостей", entities.PointTypeSpecial, false},
	}

	pointIDs := make(map[string]string)

	for i, pointData := range tavernPoints {
		pointID := fmt.Sprintf("%s_point_%d", subClusterID, i+1)
		pointIDs[pointData.name] = pointID

		point := &entities.Point{
			ID:           pointID,
			Name:         pointData.name,
			Description:  pointData.description,
			SubClusterID: subClusterID,
			Type:         pointData.pointType,
			Connections:  make([]string, 0),
			NPCs:         make([]string, 0),
			IsEntryPoint: pointData.isEntry,
		}
		subCluster.AddPoint(point)
	}

	// Настраиваем соединения согласно схеме
	// главный холл -> второй этаж / подвал / ВЫХОД
	mainHallPoint := subCluster.Points[pointIDs["Главный холл"]]
	mainHallPoint.AddConnection(pointIDs["Второй этаж"])
	mainHallPoint.AddConnection(pointIDs["Подвал"])
	if exitPointID != "" {
		mainHallPoint.AddConnection(exitPointID)
	}

	// второй этаж -> главный холл / богатые номера / бедные номера
	secondFloorPoint := subCluster.Points[pointIDs["Второй этаж"]]
	secondFloorPoint.AddConnection(pointIDs["Главный холл"])
	secondFloorPoint.AddConnection(pointIDs["Богатые номера"])
	secondFloorPoint.AddConnection(pointIDs["Бедные номера"])

	// подвал -> главный холл
	basementPoint := subCluster.Points[pointIDs["Подвал"]]
	basementPoint.AddConnection(pointIDs["Главный холл"])

	// богатые и бедные номера -> второй этаж
	richRoomsPoint := subCluster.Points[pointIDs["Богатые номера"]]
	richRoomsPoint.AddConnection(pointIDs["Второй этаж"])

	poorRoomsPoint := subCluster.Points[pointIDs["Бедные номера"]]
	poorRoomsPoint.AddConnection(pointIDs["Второй этаж"])

	// Генерируем трактирщика для главного холла
	innkeeper := tb.npcBuilder.GenerateInnkeeper(pointIDs["Главный холл"])
	world.AddNPC(innkeeper)

	// Добавляем трактирщика в Point
	if mainHallPoint, exists := subCluster.Points[pointIDs["Главный холл"]]; exists {
		mainHallPoint.NPCs = append(mainHallPoint.NPCs, innkeeper.ID)
	}

	return subCluster
}
