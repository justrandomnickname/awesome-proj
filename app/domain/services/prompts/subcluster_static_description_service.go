package prompts

import (
	"context"
	"fmt"
	"strings"

	"awesome-proj/app/domain/entities"
	"awesome-proj/app/game"
)

// SubclusterStaticDescriptionService handles AI prompt generation for subcluster static descriptions
type SubclusterStaticDescriptionService struct {
	gameEngine *game.GameEngine
}

// NewSubclusterStaticDescriptionService creates a new subcluster static description service
func NewSubclusterStaticDescriptionService(gameEngine *game.GameEngine) *SubclusterStaticDescriptionService {
	return &SubclusterStaticDescriptionService{
		gameEngine: gameEngine,
	}
}

// GenerateSubclusterDescriptionPrompt creates a structured prompt for AI content generation of subcluster descriptions
func (s *SubclusterStaticDescriptionService) GenerateSubclusterDescriptionPrompt(ctx context.Context) (string, error) {
	// Get current location hierarchy
	hierarchy, err := s.gameEngine.GetLocationHierarchy()
	if err != nil {
		return "", fmt.Errorf("не удалось получить иерархию локаций: %v", err)
	}

	// Get current point to determine context
	currentPoint, err := s.gameEngine.GetCurrentPoint()
	if err != nil {
		return "", fmt.Errorf("не удалось получить текущую точку: %v", err)
	}

	currentSubClusterID := currentPoint.SubClusterID

	// Find current subcluster and cluster
	var currentSubCluster *entities.SubCluster
	var currentCluster *entities.Cluster

	for _, cluster := range hierarchy.Clusters {
		for _, subCluster := range cluster.SubClusters {
			if subCluster.ID == currentSubClusterID {
				currentSubCluster = subCluster
				currentCluster = cluster
				break
			}
		}
		if currentSubCluster != nil {
			break
		}
	}

	if currentSubCluster == nil || currentCluster == nil {
		return "", fmt.Errorf("не удалось найти текущий субкластер")
	}

	// Collect cluster info
	clusterSubClusters := make([]string, 0, len(currentCluster.SubClusters))
	for _, subCluster := range currentCluster.SubClusters {
		clusterSubClusters = append(clusterSubClusters, subCluster.Name)
	}

	// Collect subcluster points
	subClusterPoints := make([]string, 0, len(currentSubCluster.Points))
	for _, point := range currentSubCluster.Points {
		subClusterPoints = append(subClusterPoints, point.Name)
	}

	// Collect points with connections
	pointsWithConnections := make([]map[string]interface{}, 0, len(currentSubCluster.Points))
	for _, point := range currentSubCluster.Points {
		connectedPoints := make([]string, 0)
		for _, connID := range point.Connections {
			// Find connected point name
			if connectedPoint := s.findPointByID(hierarchy, connID); connectedPoint != nil {
				if connectedPoint.Name != point.Name { // Skip self-reference
					connectedPoints = append(connectedPoints, connectedPoint.Name)
				}
			}
		}

		pointsWithConnections = append(pointsWithConnections, map[string]interface{}{
			"id":               point.ID,
			"name":             point.Name,
			"connected_points": connectedPoints,
		})
	}

	// Collect NPCs info (only real NPC data, no fallback)
	npcsWithLocations := make([]map[string]interface{}, 0)
	for _, point := range currentSubCluster.Points {
		// Only get NPCs for current point (we can only get real NPC data for current point)
		if point.ID == currentPoint.ID {
			currentNPCs, err := s.gameEngine.GetNPCsForCurrentPoint()
			if err == nil {
				for _, npc := range currentNPCs {
					npcsWithLocations = append(npcsWithLocations, map[string]interface{}{
						"id":         npc.ID,
						"name":       npc.Name,
						"race":       npc.Race,
						"point_name": point.Name,
					})
				}
			}
		}
		// For other points, we skip NPCs if we can't get real data
		// No fallback logic - if we can't get proper NPC info, we don't include them
	}

	// Generate the prompt
	prompt := s.buildSubclusterDescriptionPrompt(currentCluster, currentSubCluster, clusterSubClusters, subClusterPoints, pointsWithConnections, npcsWithLocations)

	return prompt, nil
}

// GenerateAndPrintSubclusterDescriptionPrompt creates a prompt and prints it to console
func (s *SubclusterStaticDescriptionService) GenerateAndPrintSubclusterDescriptionPrompt(ctx context.Context) (string, error) {
	prompt, err := s.GenerateSubclusterDescriptionPrompt(ctx)
	if err != nil {
		return "", err
	}

	// Print to console (terminal) only
	fmt.Println("=== ПРОМПТ ДЛЯ ГЕНЕРАЦИИ ОПИСАНИЙ СУБКЛАСТЕРА ===")
	fmt.Println(prompt)
	fmt.Println("=== КОНЕЦ ПРОМПТА ===")

	return prompt, nil
}

// findPointByID searches for a point by ID across all clusters
func (s *SubclusterStaticDescriptionService) findPointByID(hierarchy *entities.LocationHierarchy, pointID string) *entities.Point {
	for _, cluster := range hierarchy.Clusters {
		for _, subCluster := range cluster.SubClusters {
			for _, point := range subCluster.Points {
				if point.ID == pointID {
					return point
				}
			}
		}
	}
	return nil
}

// buildSubclusterDescriptionPrompt constructs the AI prompt string for subcluster descriptions
func (s *SubclusterStaticDescriptionService) buildSubclusterDescriptionPrompt(currentCluster *entities.Cluster, currentSubCluster *entities.SubCluster, clusterSubClusters []string, subClusterPoints []string, points []map[string]interface{}, npcs []map[string]interface{}) string {
	var sb strings.Builder

	sb.WriteString("Ты - опытный мастер RPG игр, специализирующийся на создании живых и детализированных игровых миров в фэнтези сеттинге.\n\n")

	sb.WriteString("ЗАДАЧА: Создай детальные описания для игрового мира на русском языке.\n\n")

	// Explanation of game structure
	sb.WriteString("СТРУКТУРА ИГРОВОГО МИРА:\n")
	sb.WriteString("- КЛАСТЕР: Большая область мира (например, деревня, город, подземелье)\n")
	sb.WriteString("- СУБКЛАСТЕР: Район внутри кластера (например, центр деревни, таверна)\n")
	sb.WriteString("- ТОЧКА (ЛОКАЦИЯ): Конкретное место в субкластере (например, зал таверны, комната)\n")
	sb.WriteString("- NPC: Персонажи, населяющие точки\n\n")

	sb.WriteString("КОНТЕКСТ КЛАСТЕРА:\n")
	sb.WriteString(fmt.Sprintf("Название кластера: \"%s\" (ID: %s)\n", currentCluster.Name, currentCluster.ID))
	sb.WriteString(fmt.Sprintf("Субкластеры в этом кластере: %s\n\n", strings.Join(clusterSubClusters, ", ")))

	sb.WriteString(fmt.Sprintf("ТЕКУЩИЙ СУБКЛАСТЕР: \"%s\" (ID: %s)\n", currentSubCluster.Name, currentSubCluster.ID))
	sb.WriteString(fmt.Sprintf("Локации в субкластере: %s\n\n", strings.Join(subClusterPoints, ", ")))

	sb.WriteString("ДЕТАЛИ ЛОКАЦИЙ:\n")
	for _, point := range points {
		pointName := point["name"].(string)
		pointID := point["id"].(string)
		connectedPoints := point["connected_points"].([]string)
		connectionsStr := "нет связей"
		if len(connectedPoints) > 0 {
			connectionsStr = strings.Join(connectedPoints, ", ")
		}
		sb.WriteString(fmt.Sprintf("- %s (ID: %s, связана с: %s)\n", pointName, pointID, connectionsStr))
	}
	sb.WriteString("\n")

	sb.WriteString("NPC В СУБКЛАСТЕРЕ:\n")
	if len(npcs) > 0 {
		for _, npc := range npcs {
			npcName := npc["name"].(string)
			npcRace := npc["race"].(string)
			npcPointName := npc["point_name"].(string)
			npcID := npc["id"].(string)
			sb.WriteString(fmt.Sprintf("- %s (%s) в локации \"%s\" (ID: %s)\n", npcName, npcRace, npcPointName, npcID))
		}
	} else {
		sb.WriteString("- В данном субкластере NPC отсутствуют или недоступны для анализа\n")
	}
	sb.WriteString("\n")

	sb.WriteString("ТРЕБОВАНИЯ К ОПИСАНИЯМ:\n")
	sb.WriteString("1. КЛАСТЕР: 3-5 предложений, общая атмосфера и суть места\n")
	sb.WriteString("2. СУБКЛАСТЕР: 3-5 предложений, более детальное описание конкретной области\n")
	sb.WriteString("3. ЛОКАЦИИ: по 3-5 предложений для каждой, уникальные детали, архитектура, атмосфера\n")
	sb.WriteString("4. NPC: по 2-3 предложения для каждого, внешность, характер, роль в мире\n\n")

	sb.WriteString("ФОРМАТ ОТВЕТА (СТРОГО СОБЛЮДАЙ):\n")
	sb.WriteString("ВАЖНО: В заголовках используй ID элементов, а НЕ их названия!\n\n")
	sb.WriteString(fmt.Sprintf("## CLUSTER_ID: %s\n", currentCluster.ID))
	sb.WriteString("[описание кластера]\n\n")

	sb.WriteString(fmt.Sprintf("## SUBCLUSTER_ID: %s\n", currentSubCluster.ID))
	sb.WriteString("[описание субкластера]\n\n")

	sb.WriteString("## LOCATIONS:\n")
	for _, point := range currentSubCluster.Points {
		sb.WriteString(fmt.Sprintf("### POINT_ID: %s\n", point.ID))
		sb.WriteString("[описание локации]\n\n")
	}

	if len(npcs) > 0 {
		sb.WriteString("## NPC:\n")
		for _, npc := range npcs {
			npcID := npc["id"].(string)
			sb.WriteString(fmt.Sprintf("### NPC_ID: %s\n", npcID))
			sb.WriteString("[описание NPC]\n\n")
		}
	}

	sb.WriteString("Создай живой, атмосферный мир с уникальными деталями для каждого элемента!")

	return sb.String()
}
