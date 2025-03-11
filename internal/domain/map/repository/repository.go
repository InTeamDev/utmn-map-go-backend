package repository

import (
	"context"
	"database/sql"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type Map struct {
	db *sql.DB
}

func NewMap(db *sql.DB) *Map {
	return &Map{db: db}
}

func (m *Map) GetObjects(ctx context.Context, req GetObjectsRequest) ([]entities.Object, error) {
	// get objects
	// get doors
	return nil, nil
}
