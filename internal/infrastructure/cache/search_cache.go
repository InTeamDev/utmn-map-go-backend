package cache

import (
	"sync"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

type InMemorySearchCache struct {
	store sync.Map
	ttl   time.Duration
}

type cacheItem struct {
	value     []entities.SearchResult
	expiresAt time.Time
}

func NewInMemorySearchCache(ttl time.Duration) *InMemorySearchCache {
	return &InMemorySearchCache{
		store: sync.Map{},
		ttl:   ttl,
	}
}

func (c *InMemorySearchCache) Set(key string, value []entities.SearchResult) {
	item := &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.store.Store(key, item)
}

func (c *InMemorySearchCache) Get(key string) ([]entities.SearchResult, bool) {
	value, ok := c.store.Load(key)
	if !ok {
		return nil, false
	}

	item, ok := value.(*cacheItem)
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		c.store.Delete(key)
		return nil, false
	}

	return item.value, true
}

func (c *InMemorySearchCache) Delete(key string) {
	c.store.Delete(key)
}

func (c *InMemorySearchCache) Clear() {
	c.store.Range(func(key, _ any) bool {
		c.store.Delete(key)
		return true
	})
}
