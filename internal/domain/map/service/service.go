package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/google/uuid"
)

//go:generate mockgen -destination=../repository/mocks/mock_map_repository.go -package=mocks github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service MapRepository
type MapRepository interface {
	GetBuildings(ctx context.Context) ([]entities.Building, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error)
	GetObjectTypes(ctx context.Context) ([]entities.ObjectType, error)
	GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (entities.GetObjectsResponse, error)
	GetObjectsByBuilding(ctx context.Context, buildingID uuid.UUID) ([]entities.Object, error)
	CreateObject(ctx context.Context, input entities.CreateObjectInput) (entities.Object, error)
	UpdateObject(ctx context.Context, input entities.UpdateObjectInput) (entities.Object, error)
	CreateBuilding(ctx context.Context, input entities.CreateBuildingInput) (entities.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
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

func (m *Map) GetObjectCategories(ctx context.Context) ([]entities.ObjectType, error) {
	objectTypes, err := m.repo.GetObjectTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object categories: %w", err)
	}
	return objectTypes, nil
}

func (m *Map) GetObjectsResponse(ctx context.Context, buildID uuid.UUID) (entities.GetObjectsResponse, error) {
	objects, err := m.repo.GetObjectsResponse(ctx, buildID)
	if err != nil {
		return entities.GetObjectsResponse{}, fmt.Errorf("get objects: %w", err)
	}
	return objects, nil
}

func (m *Map) GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]entities.Object, error) {
	objects, err := m.repo.GetObjectsByBuilding(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get objects: %w", err)
	}
	return objects, nil
}

func (m *Map) CreateObject(ctx context.Context, input entities.CreateObjectInput) (entities.Object, error) {
	object, err := m.repo.CreateObject(ctx, input)
	if err != nil {
		return entities.Object{}, fmt.Errorf("create object: %w", err)
	}
	return object, nil
}

func (m *Map) UpdateObject(ctx context.Context, input entities.UpdateObjectInput) (entities.Object, error) {
	object, err := m.repo.UpdateObject(ctx, input)
	if err != nil {
		return entities.Object{}, fmt.Errorf("get object: %w", err)
	}
	return object, nil
}

func (m *Map) CreateBuilding(ctx context.Context, input entities.CreateBuildingInput) (entities.Building, error) {
	building, err := m.repo.CreateBuilding(ctx, input)
	if err != nil {
		return entities.Building{}, fmt.Errorf("create building: %w", err)
	}
	return building, nil
}

func (m *Map) DeleteBuilding(ctx context.Context, id uuid.UUID) error {
	floors, err := m.repo.GetFloors(ctx, id)
	if err != nil {
		return fmt.Errorf("get floors: %w", err)
	}

	if len(floors) == 0 {
		return errors.New("can't delete building with floors, at first delete all floors")
	}
	err = m.repo.DeleteBuilding(ctx, id)
	if err != nil {
		return fmt.Errorf("delete building: %w", err)
	}
	return nil
}
