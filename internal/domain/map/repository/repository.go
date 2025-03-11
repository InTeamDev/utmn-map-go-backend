package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	"github.com/google/uuid"
)

type IMapConverter interface {
	ObjectsSqlcToEntity(
		objects []sqlc.GetObjectsByBuildingAndFloorRow,
		doors map[uuid.UUID][]entities.Door,
	) []entities.Object
	DoorsSqlcToEntity(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door
}

type Map struct {
	db           *sql.DB
	mapConverter IMapConverter
}

func NewMap(db *sql.DB, mapConverter IMapConverter) *Map {
	return &Map{db: db, mapConverter: mapConverter}
}

func (m *Map) GetObjects(ctx context.Context, req entities.GetObjectsRequest) ([]entities.Object, error) {
	// get objects
	q := sqlc.New(m.db)
	rowObjects, err := q.GetObjectsByBuildingAndFloor(ctx, sqlc.GetObjectsByBuildingAndFloorParams{
		BuildID: req.BuildID,
		FloorID: req.FloorID,
	})
	if err != nil {
		return nil, fmt.Errorf("get objects: %w", err)
	}
	objectIDs := make([]uuid.UUID, 0, len(rowObjects))
	for _, object := range rowObjects {
		objectIDs = append(objectIDs, object.ID)
	}
	// get doors
	rowDoors, err := q.GetDoorsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return nil, fmt.Errorf("get doors: %w", err)
	}
	// convert doors
	doors := m.mapConverter.DoorsSqlcToEntity(rowDoors)

	// convert objects
	objects := m.mapConverter.ObjectsSqlcToEntity(rowObjects, doors)
	return objects, nil
}
