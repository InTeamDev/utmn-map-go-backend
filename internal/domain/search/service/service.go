package service

import (
	"context"
	"sort"

	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/repository"
)

type SearchService struct {
	searchRepository *repository.SearchRepository
}

func NewSearchService(searchRepository *repository.SearchRepository) *SearchService {
	return &SearchService{
		searchRepository: searchRepository,
	}
}

func (s *SearchService) Search(
	ctx context.Context,
	req searchentities.SearchRequest,
) ([]searchentities.SearchResult, error) {
	searchResults, err := s.searchRepository.SearchObjectsWithDoors(ctx, req.BuildingID, req.Query)
	if err != nil {
		return nil, err
	}

	// Сортируем по Preview (а при равенстве — по UUID)
	sort.Slice(searchResults, func(i, j int) bool {
		if searchResults[i].Preview != searchResults[j].Preview {
			return searchResults[i].Preview < searchResults[j].Preview
		}
		return searchResults[i].ObjectID.String() < searchResults[j].ObjectID.String()
	})

	return searchResults, nil
}
