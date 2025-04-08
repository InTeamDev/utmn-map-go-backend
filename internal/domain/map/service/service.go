package service

import (
	"context"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/google/uuid"
)

type MapRepository interface {
	GetBuildings(ctx context.Context) ([]entities.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error)
	GetObjectTypes(ctx context.Context, buildID uuid.UUID) ([]entities.ObjectType, error)
	GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]entities.Object, error)
}

type Map struct {
	repo MapRepository
}

func NewMap(repo MapRepository) *Map {
	return &Map{repo: repo}
}

func (m *Map) GetBuildings(ctx context.Context) ([]entities.Building, error) {
	buildings, err := m.repo.GetBuildings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}
	return buildings, nil
}

func (m *Map) GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error) {
	floors, err := m.repo.GetFloors(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get floors: %w", err)
	}
	return floors, nil
}

func (m *Map) GetObjectCategories(ctx context.Context, buildID uuid.UUID) ([]entities.ObjectType, error) {
	objectTypes, err := m.repo.GetObjectTypes(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get object categories: %w", err)
	}
	return objectTypes, nil
}

func (m *Map) GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]entities.Object, error) {
	objects, err := m.repo.GetObjectsByBuilding(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get objects: %w", err)
	}
	return objects, nil
}
