package services

import (
	"math/rand"

	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/infrastructure/builders"
)

// WorldGenerationService handles world generation logic
type WorldGenerationService struct {
	locationBuilder *builders.LocationBuilder
}

// NewWorldGenerationService creates a new world generation service
func NewWorldGenerationService() *WorldGenerationService {
	return &WorldGenerationService{
		locationBuilder: builders.NewLocationBuilder(),
	}
}

// GenerateWorld creates a new world with random content
func (wgs *WorldGenerationService) GenerateWorld(name string, seed int64) *aggregates.World {
	world := aggregates.NewEmptyWorld(name, seed)
	
	// Generate random number generator with seed
	rng := rand.New(rand.NewSource(seed))
	
	// Generate locations using infrastructure builder
	wgs.locationBuilder.GenerateRandomLocations(world, rng, 3)
	
	return world
}
