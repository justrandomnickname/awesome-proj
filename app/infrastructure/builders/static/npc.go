package static

import (
	"awesome-proj/app/domain/entities"
	"fmt"
)

func GenerateInnkeeper(pointID string) *entities.NPC {
	return &entities.NPC{
		ID:          fmt.Sprintf("%s_innkeeper", pointID),
		Name:        "Трактирщик",
		Race:        "human",
		LocationID:  pointID,
		Description: "NOT_SPECIFIED",
	}
}
