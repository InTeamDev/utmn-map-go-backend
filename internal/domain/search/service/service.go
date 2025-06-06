package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	mapcache "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/cache"
	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	mapservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

type SearchService struct {
	mapCache   *mapcache.InMemoryMapCache
	mapService *mapservice.Map
}

func NewSearchService(mapCache *mapcache.InMemoryMapCache, mapService *mapservice.Map) *SearchService {
	return &SearchService{
		mapCache:   mapCache,
		mapService: mapService,
	}
}

func (s *SearchService) Search(
	ctx context.Context,
	req searchentities.SearchRequest,
) ([]searchentities.SearchResult, error) {
	objects, exists := s.mapCache.Get(req.BuildingID)
	if !exists {
		var err error
		objects, err = s.mapService.GetObjectsByBuilding(ctx, req.BuildingID)
		if err != nil {
			return nil, fmt.Errorf("get objects: %w", err)
		}
		s.mapCache.Set(req.BuildingID, objects)
	}

	objectDoorMap, err := s.mapService.GetObjectDoorPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object door pairs: %w", err)
	}

	query := strings.ToLower(req.Query)

	filtered := make([]searchentities.SearchResult, 0, len(objects))
	for _, obj := range objects {
		doorID, hasDoor := objectDoorMap[obj.ID]
		if !hasDoor {
			continue
		}
		if matchesQuery(obj, query) {
			filtered = append(filtered, searchentities.SearchResult{
				ObjectID:     obj.ID,
				ObjectTypeID: int(obj.ObjectTypeID),
				Preview:      fmt.Sprintf("%s (%s)", obj.Name, obj.Floor.Name),
				DoorID:       doorID,
			})
		}
	}

	// Сортируем по Preview (а при равенстве — по UUID)
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Preview != filtered[j].Preview {
			return filtered[i].Preview < filtered[j].Preview
		}
		return filtered[i].ObjectID.String() < filtered[j].ObjectID.String()
	})

	return filtered, nil
}

// matchesQuery проверяет, содержит ли объект obj строку query
// (по имени, по алиасу или по описанию). Если query == "", то всегда true.
func matchesQuery(obj mapentities.Object, query string) bool {
	if query == "" {
		return true
	}
	nameLower := strings.ToLower(obj.Name)
	aliasLower := strings.ToLower(obj.Alias)
	descLower := strings.ToLower(obj.Description)

	return strings.Contains(nameLower, query) ||
		strings.Contains(aliasLower, query) ||
		strings.Contains(descLower, query)
}
