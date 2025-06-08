package service

import (
	"context"
	"sort"
	"time"

	"github.com/patrickmn/go-cache"

	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/repository"
)

type SearchService struct {
	searchRepository *repository.SearchRepository
	cache            *cache.Cache
}

func NewSearchService(searchRepository *repository.SearchRepository) *SearchService {
	return &SearchService{
		searchRepository: searchRepository,
		cache:            cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (s *SearchService) Search(
	ctx context.Context,
	req searchentities.SearchRequest,
) ([]searchentities.SearchResult, error) {
	key := req.BuildingID.String() + "|" + req.Query

	if cached, found := s.cache.Get(key); found {
		if results, ok := cached.([]searchentities.SearchResult); ok {
			return results, nil
		}
		s.cache.Delete(key)
	}

	searchResults, err := s.searchRepository.SearchObjectsWithDoors(ctx, req.BuildingID, req.Query)
	if err != nil {
		return nil, err
	}

	s.cache.Set(key, searchResults, cache.DefaultExpiration)

	sort.Slice(searchResults, func(i, j int) bool {
		if searchResults[i].Preview != searchResults[j].Preview {
			return searchResults[i].Preview < searchResults[j].Preview
		}
		return searchResults[i].ObjectID.String() < searchResults[j].ObjectID.String()
	})
	return searchResults, nil
}
