package services

import (
	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/infrastructure/builders"
	"math/rand"
)

type WorldGenerationService struct {
	locationBuilder *builders.LocationBuilder
}

func NewWorldGenerationService() *WorldGenerationService {
	return &WorldGenerationService{
		locationBuilder: builders.NewLocationBuilder(),
	}
}
func (wgs *WorldGenerationService) GenerateWorld(name string, seed int64) *aggregates.World {
	world := aggregates.NewEmptyWorld(name, seed)
	rng := rand.New(rand.NewSource(seed))

	hierarchy := wgs.locationBuilder.GenerateLocationHierarchyWithNPCs(world, rng)
	world.SetHierarchy(hierarchy)

	return world
}
