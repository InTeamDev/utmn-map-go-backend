package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	"github.com/google/uuid"
)

type MapConverter interface {
	ObjectsSqlcToEntityByBuilding(
		objects []sqlc.GetObjectsByBuildingRow,
		doors map[uuid.UUID][]entities.Door,
	) []entities.Object
	DoorsSqlcToEntity(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door
	FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor
	BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building
	ObjectTypesSqlcToEntity(objectTypes []sqlc.ObjectType) []entities.ObjectType
}

type Map struct {
	db           *sql.DB
	mapConverter MapConverter
}

func NewMap(db *sql.DB, mapConverter MapConverter) *Map {
	return &Map{db: db, mapConverter: mapConverter}
}

func (m *Map) GetBuildings(ctx context.Context) ([]entities.Building, error) {
	q := sqlc.New(m.db)
	buildings, err := q.GetBuildings(ctx)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}
	return m.mapConverter.BuildingsSqlcToEntity(buildings), nil
}

func (m *Map) GetFloors(ctx context.Context, buildID uuid.UUID) ([]entities.Floor, error) {
	q := sqlc.New(m.db)
	floors, err := q.GetFloorsByBuilding(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get floors: %w", err)
	}
	return m.mapConverter.FloorsSqlcToEntity(floors), nil
}

func (m *Map) GetObjectTypes(ctx context.Context, buildID uuid.UUID) ([]entities.ObjectType, error) {
	q := sqlc.New(m.db)
	objectTypes, err := q.GetObjectTypes(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get object types: %w", err)
	}
	return m.mapConverter.ObjectTypesSqlcToEntity(objectTypes), nil
}

func (m *Map) GetObjectsByBuilding(ctx context.Context, buildID uuid.UUID) ([]entities.Object, error) {
	q := sqlc.New(m.db)
	rowObjects, err := q.GetObjectsByBuilding(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("get objects by building: %w", err)
	}
	objectIDs := make([]uuid.UUID, 0, len(rowObjects))
	for _, object := range rowObjects {
		objectIDs = append(objectIDs, object.ID)
	}
	rowDoors, err := q.GetDoorsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return nil, fmt.Errorf("get doors: %w", err)
	}
	doors := m.mapConverter.DoorsSqlcToEntity(rowDoors)
	objects := m.mapConverter.ObjectsSqlcToEntityByBuilding(rowObjects, doors)
	return objects, nil
}
