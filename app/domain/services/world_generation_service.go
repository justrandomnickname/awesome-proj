package services
import (
	"math/rand"
	"awesome-proj/app/domain/aggregates"
	"awesome-proj/app/infrastructure/builders"
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
	wgs.locationBuilder.GenerateRandomLocations(world, rng, 1)
	return world
}
