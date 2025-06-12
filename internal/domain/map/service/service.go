package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

//go:generate mockgen -destination=../repository/mocks/mock_map_repository.go -package=mocks github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service MapRepository
type MapRepository interface {
	GetObjectByID(ctx context.Context, objectID uuid.UUID) (entities.Object, error)
	GetBuildings(ctx context.Context) ([]entities.Building, error)
	GetPolygonByID(ctx context.Context, id uuid.UUID) (entities.Polygon, error)
	GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error)
	GetDoors(ctx context.Context, buildID uuid.UUID) ([]entities.GetDoorsResponse, error)
	GetObjectTypes(ctx context.Context) ([]entities.ObjectTypeInfo, error)
	GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (entities.GetObjectsResponse, error)
	GetObjectsByBuilding(ctx context.Context, buildingID uuid.UUID) ([]entities.Object, error)
	GetObjectTypeByID(ctx context.Context, id int32) (entities.ObjectTypeInfo, error)
	CreateObject(
		ctx context.Context,
		floorID uuid.UUID,
		input entities.CreateObjectInput,
	) (entities.Object, error)
	UpdateObject(ctx context.Context, id uuid.UUID, input entities.UpdateObjectInput) (entities.Object, error)
	DeleteObject(ctx context.Context, objectID uuid.UUID) error
	GetDoor(
		ctx context.Context,
		buildingID uuid.UUID,
		doorID uuid.UUID,
	) (entities.Door, error)
	UpdateDoor(
		ctx context.Context,
		buildingID uuid.UUID,
		doorID uuid.UUID,
		input entities.Door,
	) (entities.Door, error)
	CreateBuilding(ctx context.Context, input entities.CreateBuildingInput) (entities.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID) error
	UpdateBuilding(ctx context.Context, id uuid.UUID, input entities.UpdateBuildingInput) (entities.Building, error)
	GetBuildingByID(ctx context.Context, id uuid.UUID) (entities.Building, error)
	CreateFloor(ctx context.Context, buildingID uuid.UUID, floor entities.Floor) error
	CreateDoor(ctx context.Context, objectID uuid.UUID, door entities.Door) (entities.Door, error)
	CreatePolygon(ctx context.Context, polygon entities.Polygon) (entities.Polygon, error)
	CreatePolygonPoint(
		ctx context.Context,
		polygonID uuid.UUID,
		order int32,
		x, y float64,
	) (entities.PolygonPoint, error)
	GetDoorFloorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	GetObjectDoorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error)
	DeletePolygonPoints(ctx context.Context, request entities.DeletePolygonPointsRequest) error
	GetPolygonsByFloorID(ctx context.Context, floorID uuid.UUID) ([]entities.Polygon, error)
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

func (m *Map) GetDoors(ctx context.Context, buildID uuid.UUID) ([]entities.GetDoorsResponse, error) {
	doors, err := m.repo.GetDoors(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get doors: %w", err)
	}
	return doors, nil
}

func (m *Map) GetDoor(
	ctx context.Context,
	buildingID uuid.UUID,
	doorID uuid.UUID,
) (entities.Door, error) {
	door, err := m.repo.GetDoor(ctx, buildingID, doorID)
	if err != nil {
		return entities.Door{}, fmt.Errorf("get door: %w", err)
	}
	return door, nil
}

func (m *Map) UpdateDoor(
	ctx context.Context,
	buildingID uuid.UUID,
	doorID uuid.UUID,
	input entities.Door,
) (entities.Door, error) {
	if input.Width <= 0 || input.Height <= 0 {
		return entities.Door{}, fmt.Errorf("%w: width and height must be positive", entities.ErrInvalidInput)
	}

	if input.X < 0 || input.Y < 0 {
		return entities.Door{}, fmt.Errorf("%w: coordinates cannot be negative", entities.ErrInvalidInput)
	}

	if _, err := m.repo.GetDoor(ctx, buildingID, doorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Door{}, entities.ErrDoorNotFound
		}
		return entities.Door{}, err
	}

	return m.repo.UpdateDoor(ctx, buildingID, doorID, input)
}

func (m *Map) GetObjectCategories(ctx context.Context) ([]entities.ObjectTypeInfo, error) {
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

func (m *Map) GetObjectTypeByID(ctx context.Context, id int32) (entities.ObjectTypeInfo, error) {
	inputID := id

	objectType, err := m.repo.GetObjectTypeByID(ctx, inputID)
	if err != nil {
		return entities.ObjectTypeInfo{}, fmt.Errorf("get object type by id: %w", err)
	}
	return objectType, nil
}

func (m *Map) GetObjectByID(
	ctx context.Context,
	objectID uuid.UUID,
) (entities.Object, error) {
	object, err := m.repo.GetObjectByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Object{}, entities.ErrObjectNotFound
		}
		return entities.Object{}, fmt.Errorf("get object by id: %w", err)
	}

	return object, nil
}

func (m *Map) CreateObject(
	ctx context.Context,
	floorID uuid.UUID,
	input entities.CreateObjectInput,
) (entities.Object, error) {
	_, err := m.GetObjectTypeByID(ctx, input.ObjectTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Object{}, entities.ErrObjectTypeNotFound
		}
		return entities.Object{}, fmt.Errorf("get object type: %w", err)
	}

	object, err := m.repo.CreateObject(ctx, floorID, input)
	if err != nil {
		return entities.Object{}, fmt.Errorf("create object: %w", err)
	}

	return object, nil
}

func (m *Map) UpdateObject(
	ctx context.Context,
	id uuid.UUID,
	input entities.UpdateObjectInput,
) (entities.Object, error) {
	if input.ObjectTypeID != nil {
		objectType, err := m.GetObjectTypeByID(ctx, *input.ObjectTypeID)
		if err != nil {
			return entities.Object{}, fmt.Errorf("object type validation failed: %w", err)
		}
		input.ObjectTypeID = &objectType.ID
	}

	object, err := m.repo.UpdateObject(ctx, id, input)
	if err != nil {
		return entities.Object{}, fmt.Errorf("update object: %w", err)
	}
	return object, nil
}

func (m *Map) DeleteObject(ctx context.Context, objectID uuid.UUID) error {
	err := m.repo.DeleteObject(ctx, objectID)
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
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

	if len(floors) != 0 {
		return errors.New("can't delete building with floors, at first delete all floors")
	}

	err = m.repo.DeleteBuilding(ctx, id)
	if err != nil {
		return fmt.Errorf("delete building: %w", err)
	}
	return nil
}

func (m *Map) UpdateBuilding(
	ctx context.Context,
	id uuid.UUID,
	input entities.UpdateBuildingInput,
) (entities.Building, error) {
	building, err := m.repo.UpdateBuilding(ctx, id, input)
	if err != nil {
		return entities.Building{}, fmt.Errorf("get building: %w", err)
	}
	return building, nil
}

func (m *Map) GetBuildingByID(ctx context.Context, id uuid.UUID) (entities.Building, error) {
	return m.repo.GetBuildingByID(ctx, id)
}

func (m *Map) CreatePolygon(ctx context.Context, polygon entities.Polygon) (entities.Polygon, error) {
	return m.repo.CreatePolygon(ctx, polygon)
}

func (m *Map) CreatePolygonPoint(
	ctx context.Context,
	polygonID uuid.UUID,
	order int32,
	x, y float64,
) (entities.PolygonPoint, error) {
	return m.repo.CreatePolygonPoint(ctx, polygonID, order, x, y)
}

func (m *Map) CreateFloor(ctx context.Context, buildingID uuid.UUID, floor entities.Floor) error {
	return m.repo.CreateFloor(ctx, buildingID, floor)
}

func (m *Map) CreateDoor(ctx context.Context, objectID uuid.UUID, door entities.Door) (entities.Door, error) {
	door, err := m.repo.CreateDoor(ctx, objectID, door)
	if err != nil {
		return entities.Door{}, fmt.Errorf("create door: %w", err)
	}

	return door, nil
}

func (m *Map) GetDoorFloorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	doorFloorPairs, err := m.repo.GetDoorFloorPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get door-floor pairs: %w", err)
	}
	return doorFloorPairs, nil
}

func (m *Map) GetObjectDoorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	objectDoorPairs, err := m.repo.GetObjectDoorPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object-door pairs: %w", err)
	}
	return objectDoorPairs, nil
}

func (m *Map) DeletePolygonPoints(ctx context.Context, request entities.DeletePolygonPointsRequest) error {
	return m.repo.DeletePolygonPoints(ctx, request)
}

func (m *Map) GetPolygonByID(ctx context.Context, id uuid.UUID) (entities.Polygon, error) {
	polygon, err := m.repo.GetPolygonByID(ctx, id)
	if err != nil {
		return entities.Polygon{}, err
	}

	return polygon, nil
}

func (m *Map) GetPolygonsByFloorID(ctx context.Context, floorID uuid.UUID) ([]entities.Polygon, error) {
	return m.repo.GetPolygonsByFloorID(ctx, floorID)
}
