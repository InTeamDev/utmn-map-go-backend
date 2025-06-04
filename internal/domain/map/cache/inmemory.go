package cache

import (
	"github.com/google/uuid"

	mapentites "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type InMemoryMapCache struct {
	store map[uuid.UUID][]mapentites.Object
}

func NewInMemoryMapCache() *InMemoryMapCache {
	return &InMemoryMapCache{
		store: make(map[uuid.UUID][]mapentites.Object),
	}
}

func (c *InMemoryMapCache) Set(buildID uuid.UUID, objects []mapentites.Object) {
	c.store[buildID] = objects
}

func (c *InMemoryMapCache) Get(buildID uuid.UUID) ([]mapentites.Object, bool) {
	objects, ok := c.store[buildID]
	return objects, ok
}
