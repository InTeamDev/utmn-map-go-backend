package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/repository/sqlc"
)

type SearchRepository struct {
	q  *sqlc.Queries
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		q:  sqlc.New(db),
		db: db,
	}
}

func (s *SearchRepository) SearchObjectsWithDoors(
	ctx context.Context,
	buildingID uuid.UUID,
	query string,
) ([]entities.SearchResult, error) {
	objects, err := s.q.Search(ctx, sqlc.SearchParams{
		ID:      buildingID,
		Column2: query,
	})
	if err != nil {
		return nil, err
	}
	results := make([]entities.SearchResult, 0, len(objects))
	for _, obj := range objects {
		preview, ok := obj.Preview.(string)
		if !ok {
			preview = fmt.Sprint(obj.Preview)
		}
		results = append(results, entities.SearchResult{
			ObjectID:     obj.ObjectID,
			ObjectTypeID: int(obj.ObjectTypeID),
			Preview:      preview,
			DoorID:       obj.DoorID,
		})
	}
	return results, nil
}
