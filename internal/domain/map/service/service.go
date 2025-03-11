package service

import (
	"context"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type MapRepository interface {
	GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error)
}

type Map struct {
	repo MapRepository
}

func NewMap(repo MapRepository) *Map {
	return &Map{repo: repo}
}

func (m *Map) GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error) {
	objects, err := m.repo.GetObjects(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get objects: %w", err)
	}
	return objects, nil
}
