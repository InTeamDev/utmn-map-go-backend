package service

import (
	"context"
	"log"
	"sort"
	"strings"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	maprepository "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/popularity/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/utils"
	"github.com/InTeamDev/utmn-map-go-backend/internal/infrastructure/cache"
	"github.com/google/uuid"
)

const RANK = 0.1

type SearchService struct {
	cache            *cache.InMemorySearchCache
	repo             *maprepository.Map
	queryProcessor   *utils.QueryProcessor
	popularityRanker *service.PopularityRanker
}

func NewSearchService(
	cache *cache.InMemorySearchCache,
	repo *maprepository.Map,
	queryProcessor *utils.QueryProcessor,
	popularityRanker *service.PopularityRanker,
) *SearchService {
	return &SearchService{
		cache:            cache,
		repo:             repo,
		queryProcessor:   queryProcessor,
		popularityRanker: popularityRanker,
	}
}

func (s *SearchService) Search(
	ctx context.Context,
	query string,
	floorID uuid.UUID,
) ([]searchentities.SearchResult, error) {
	cacheKey := utils.GenerateCacheKey(query, floorID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached, nil
	}

	processedQuery := s.queryProcessor.Process(query)

	objects, err := s.repo.GetObjectsByFloor(ctx, floorID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []searchentities.SearchResult
	for _, obj := range objects {
		if s.isMatch(processedQuery, obj) && s.isSameFloor(obj, floorID) {
			results = append(results, s.buildResult(obj))
		}
	}

	s.sortResults(results, floorID)
	s.cache.Set(cacheKey, results)
	return results, nil
}

func (s *SearchService) isMatch(query string, obj mapentities.Object) bool {
	return strings.Contains(strings.ToLower(obj.Name), query)
}

func (s *SearchService) buildResult(
	obj mapentities.Object,
) searchentities.SearchResult {
	result := searchentities.SearchResult{
		ID:         obj.ID,
		Popularity: s.popularityRanker.GetPopularityScore(obj.ID),
		FloorID:    obj.Floor.ID,
		Type:       string(obj.ObjectType),
		X:          obj.X,
		Y:          obj.Y,
	}

	return result
}

func (s *SearchService) isSameFloor(obj mapentities.Object, userFloor uuid.UUID) bool {
	return obj.Floor.ID == userFloor
}

func (s *SearchService) sortResults(
	results []searchentities.SearchResult,
	userFloor uuid.UUID,
) {
	sort.Slice(results, func(i, j int) bool {
		a, b := results[i], results[j]

		if a.FloorID != b.FloorID {
			return a.FloorID == userFloor
		}

		if a.Relevance != b.Relevance {
			return a.Relevance > b.Relevance
		}

		if a.Popularity != b.Popularity {
			return a.Popularity > b.Popularity
		}

		return a.Distance < b.Distance
	})
}
