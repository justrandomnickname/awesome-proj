package static

import (
	"awesome-proj/app/domain/entities/npc"
	"fmt"
)

func GenerateInnkeeper(pointID string) *npc.NPC {
	return &npc.NPC{
		ID:          fmt.Sprintf("%s_innkeeper", pointID),
		Name:        "Трактирщик",
		Race:        "human",
		LocationID:  pointID,
		Description: "NOT_SPECIFIED",
		TemperTraits: npc.TemperTraits{
			ID:             "",
			NPC_ID:         "",
			Prudence:       0,
			Emotionality:   0,
			Independence:   0,
			Optimism:       0,
			Flexibility:    0,
			Aggressiveness: 0,
		},
		TraitIDs: nil,
	}
}
