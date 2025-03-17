package service

import (
	"sort"
	"strings"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/popularity/service"
	repository "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/repository/search"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/utils"
	"github.com/InTeamDev/utmn-map-go-backend/internal/infrastructure/cache"
)

type SearchService struct {
	cache            *cache.InMemorySearchCache
	repo             repository.SearchRepository
	queryProcessor   *utils.QueryProcessor
	relevanceCalc    *utils.RelevanceCalculator
	distanceService  *utils.DistanceService
	popularityRanker *service.PopularityRanker
}

func (s *SearchService) Search(
	query string,
	userFloor string,
	ctx *searchentities.UserContext,
) ([]searchentities.SearchResult, error) {
	cacheKey := utils.GenerateCacheKey(query, userFloor, ctx)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached, nil
	}

	processedQuery := s.queryProcessor.Process(query)
	objects, _ := s.repo.GetAllObjects()

	var results []searchentities.SearchResult
	for _, obj := range objects {
		if s.isMatch(processedQuery, obj) && s.isSameFloor(obj, userFloor) {
			relevance := s.relevanceCalc.Calculate(processedQuery, obj, ctx)
			if relevance > 0.5 {
				results = append(results, s.buildResult(obj, relevance, ctx))
			}
		}
	}

	s.sortResults(results, userFloor, ctx)
	s.cache.Set(cacheKey, results)
	return results, nil
}

// isMatch проверяет, соответствует ли объект поисковому запросу.
func (s *SearchService) isMatch(query string, obj mapentities.Object) bool {
	return strings.Contains(strings.ToLower(string(obj.ObjectType)), query) ||
		strings.Contains(strings.ToLower(obj.Floor.Name), query)
}

// buildResult создает SearchResult на основе объекта и релевантности.
func (s *SearchService) buildResult(
	obj mapentities.Object,
	relevance float64,
	ctx *searchentities.UserContext,
) searchentities.SearchResult {
	result := searchentities.SearchResult{
		ID:         obj.ID,
		Relevance:  relevance,
		Popularity: s.popularityRanker.GetPopularityScore(obj.ID),
		Floor:      obj.Floor.Name,
		Type:       string(obj.ObjectType),
		X:          obj.X,
		Y:          obj.Y,
	}

	if ctx != nil && ctx.Location != nil {
		result.Distance = s.distanceService.Calculate(*ctx.Location, obj)
	}

	return result
}

func (s *SearchService) isSameFloor(obj mapentities.Object, userFloor string) bool {
	return strings.EqualFold(obj.Floor.Name, userFloor)
}

func (s *SearchService) sortResults(
	results []searchentities.SearchResult,
	userFloor string,
	ctx *searchentities.UserContext,
) {
	sort.Slice(results, func(i, j int) bool {
		a, b := results[i], results[j]

		if a.Floor != b.Floor {
			return a.Floor == userFloor
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
