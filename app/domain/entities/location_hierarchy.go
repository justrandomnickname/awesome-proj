package entities

type LocationHierarchy struct {
	Clusters map[string]*Cluster `json:"clusters"`
}

type Cluster struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	MainPoint     string                 `json:"main_point"`
	SubClusters   map[string]*SubCluster `json:"sub_clusters"`
	ChildClusters map[string]*Cluster    `json:"child_clusters,omitempty"`
}

type SubCluster struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ClusterID   string            `json:"cluster_id"`
	EntryPoints []string          `json:"entry_points"`
	Points      map[string]*Point `json:"points"`
}

type Point struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	SubClusterID string    `json:"sub_cluster_id"`
	Type         PointType `json:"type"`
	Connections  []string  `json:"connections"`
	NPCs         []string  `json:"npcs"`
	IsEntryPoint bool      `json:"is_entry_point"`
}

type ConnectionInfo struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	SubClusterID     string            `json:"sub_cluster_id"`
	Type             PointType         `json:"type"`
	Connections      []string          `json:"connections"`
	ConnectionNames  map[string]string `json:"connection_names"` // ID -> Name соединений
	NPCs             []string          `json:"npcs"`
	IsEntryPoint     bool              `json:"is_entry_point"`
	DisplayName      string            `json:"display_name"`       // Отображаемое имя (может быть "Перейти в Таверна")
	IsInterCluster   bool              `json:"is_inter_cluster"`   // Переход между субкластерами
	TargetSubCluster string            `json:"target_sub_cluster"` // Название целевого субкластера
}

type PointType string

const (
	PointTypeEntry   PointType = "entry"
	PointTypeRegular PointType = "regular"
	PointTypeSpecial PointType = "special"
	PointTypeExit    PointType = "exit"
)

func NewLocationHierarchy() *LocationHierarchy {
	return &LocationHierarchy{
		Clusters: make(map[string]*Cluster),
	}
}

func (lh *LocationHierarchy) AddCluster(cluster *Cluster) {
	lh.Clusters[cluster.ID] = cluster
}

func (c *Cluster) AddSubCluster(subCluster *SubCluster) {
	if c.SubClusters == nil {
		c.SubClusters = make(map[string]*SubCluster)
	}
	c.SubClusters[subCluster.ID] = subCluster
}

func (sc *SubCluster) AddPoint(point *Point) {
	if sc.Points == nil {
		sc.Points = make(map[string]*Point)
	}
	sc.Points[point.ID] = point

	if point.IsEntryPoint {
		sc.EntryPoints = append(sc.EntryPoints, point.ID)
	}
}

func (p *Point) AddConnection(pointID string) {
	for _, existing := range p.Connections {
		if existing == pointID {
			return
		}
	}
	p.Connections = append(p.Connections, pointID)
}
func (lh *LocationHierarchy) FindPoint(pointID string) *Point {
	for _, cluster := range lh.Clusters {
		for _, subCluster := range cluster.SubClusters {
			if point, exists := subCluster.Points[pointID]; exists {
				return point
			}
		}
	}
	return nil
}
func (lh *LocationHierarchy) FindCluster(clusterID string) *Cluster {
	return lh.Clusters[clusterID]
}
func (lh *LocationHierarchy) FindSubCluster(subClusterID string) *SubCluster {
	for _, cluster := range lh.Clusters {
		if subCluster, exists := cluster.SubClusters[subClusterID]; exists {
			return subCluster
		}
	}
	return nil
}

func (lh *LocationHierarchy) FindSubClusterByPoint(pointID string) *SubCluster {
	for _, cluster := range lh.Clusters {
		for _, subCluster := range cluster.SubClusters {
			if _, exists := subCluster.Points[pointID]; exists {
				return subCluster
			}
		}
	}
	return nil
}

// NewLocationHierarchyFromMap создает LocationHierarchy из map[string]interface{}
func NewLocationHierarchyFromMap(hierarchyData interface{}) *LocationHierarchy {
	hierarchyMap, ok := hierarchyData.(map[string]interface{})
	if !ok {
		return nil
	}

	hierarchy := NewLocationHierarchy()

	// Загружаем clusters
	if clustersData, ok := hierarchyMap["clusters"].(map[string]interface{}); ok {
		for clusterID, clusterData := range clustersData {
			if cluster := newClusterFromMap(clusterData); cluster != nil {
				cluster.ID = clusterID
				hierarchy.AddCluster(cluster)
			}
		}
	}

	return hierarchy
}

// newClusterFromMap создает Cluster из map
func newClusterFromMap(clusterData interface{}) *Cluster {
	clusterMap, ok := clusterData.(map[string]interface{})
	if !ok {
		return nil
	}

	cluster := &Cluster{
		ID:            getStringFromMapH(clusterMap, "id"),
		Name:          getStringFromMapH(clusterMap, "name"),
		Description:   getStringFromMapH(clusterMap, "description"),
		Type:          getStringFromMapH(clusterMap, "type"),
		MainPoint:     getStringFromMapH(clusterMap, "main_point"),
		SubClusters:   make(map[string]*SubCluster),
		ChildClusters: make(map[string]*Cluster),
	}

	// Загружаем sub_clusters
	if subClustersData, ok := clusterMap["sub_clusters"].(map[string]interface{}); ok {
		for subClusterID, subClusterData := range subClustersData {
			if subCluster := newSubClusterFromMap(subClusterData); subCluster != nil {
				subCluster.ID = subClusterID
				cluster.AddSubCluster(subCluster)
			}
		}
	}

	return cluster
}

// newSubClusterFromMap создает SubCluster из map
func newSubClusterFromMap(subClusterData interface{}) *SubCluster {
	subClusterMap, ok := subClusterData.(map[string]interface{})
	if !ok {
		return nil
	}

	subCluster := &SubCluster{
		ID:          getStringFromMapH(subClusterMap, "id"),
		Name:        getStringFromMapH(subClusterMap, "name"),
		Description: getStringFromMapH(subClusterMap, "description"),
		ClusterID:   getStringFromMapH(subClusterMap, "cluster_id"),
		EntryPoints: getStringArrayFromMap(subClusterMap, "entry_points"),
		Points:      make(map[string]*Point),
	}

	// Загружаем points
	if pointsData, ok := subClusterMap["points"].(map[string]interface{}); ok {
		for pointID, pointData := range pointsData {
			if point := newPointFromMap(pointData); point != nil {
				point.ID = pointID
				subCluster.AddPoint(point)
			}
		}
	}

	return subCluster
}

// newPointFromMap создает Point из map
func newPointFromMap(pointData interface{}) *Point {
	pointMap, ok := pointData.(map[string]interface{})
	if !ok {
		return nil
	}

	return &Point{
		ID:           getStringFromMapH(pointMap, "id"),
		Name:         getStringFromMapH(pointMap, "name"),
		Description:  getStringFromMapH(pointMap, "description"),
		SubClusterID: getStringFromMapH(pointMap, "sub_cluster_id"),
		Type:         PointType(getStringFromMapH(pointMap, "type")),
		Connections:  getStringArrayFromMap(pointMap, "connections"),
		NPCs:         getStringArrayFromMap(pointMap, "npcs"),
		IsEntryPoint: getBoolFromMap(pointMap, "is_entry_point"),
	}
}

// Вспомогательные функции для извлечения данных из map
func getStringFromMapH(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getStringArrayFromMap(m map[string]interface{}, key string) []string {
	if arr, ok := m[key].([]interface{}); ok {
		result := make([]string, len(arr))
		for i, v := range arr {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}
	return []string{}
}

func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}
