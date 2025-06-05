package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	mapcache "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/cache"
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
	// simple search
	objects, exists := s.mapCache.Get(req.BuildingID)
	if !exists {
		var err error
		objects, err = s.mapService.GetObjectsByBuilding(ctx, req.BuildingID)
		if err != nil {
			return nil, fmt.Errorf("get objects: %w", err)
		}
		s.mapCache.Set(req.BuildingID, objects)
	}

	// filter objects by name
	filteredObjects := make([]searchentities.SearchResult, 0, len(objects))
	for _, object := range objects {
		if strings.Contains(strings.ToLower(object.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(object.Alias), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(object.Description), strings.ToLower(req.Query)) ||
			req.Query == "" {
			filteredObjects = append(filteredObjects, searchentities.SearchResult{
				ObjectID:     object.ID,
				ObjectTypeID: int(object.ObjectTypeID),
				Preview:      fmt.Sprintf("%s (%s)", object.Name, object.Floor.Name),
			})
		}
	}

	sort.Slice(filteredObjects, func(i, j int) bool {
		// Sort by name, then by alias, then by description
		if filteredObjects[i].Preview != filteredObjects[j].Preview {
			return filteredObjects[i].Preview < filteredObjects[j].Preview
		}
		return filteredObjects[i].ObjectID.String() < filteredObjects[j].ObjectID.String()
	})

	return filteredObjects, nil
}
