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
	ObjectSqlcToEntity(object sqlc.GetObjectsByBuildingRow, doors []entities.Door) entities.Object
	DoorsSqlcToEntityMap(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door
	FloorSqlcToEntity(f sqlc.Floor) entities.Floor
	FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor
	BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building
	ObjectTypesSqlcToEntity(objectTypes []sqlc.ObjectType) []entities.ObjectType
	// Новая функция для конвертации background этажа
	FloorBackgroundSqlcToEntityMany(rows []sqlc.GetFloorBackgroundRow) []entities.FloorBackgroundElement
}

type Map struct {
	q         *sqlc.Queries
	converter MapConverter
}

func NewMap(db *sql.DB, converter MapConverter) *Map {
	return &Map{
		q:         sqlc.New(db),
		converter: converter,
	}
}

func (r *Map) GetBuildings(ctx context.Context) ([]entities.Building, error) {
	buildings, err := r.q.GetBuildings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get buildings: %w", err)
	}
	return r.converter.BuildingsSqlcToEntity(buildings), nil
}

func (r *Map) GetFloors(ctx context.Context, buildingID uuid.UUID) ([]entities.Floor, error) {
	floors, err := r.q.GetFloorsByBuilding(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get floors: %w", err)
	}
	return r.converter.FloorsSqlcToEntity(floors), nil
}

func (r *Map) GetObjectTypes(ctx context.Context) ([]entities.ObjectType, error) {
	objectTypes, err := r.q.GetObjectTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object types: %w", err)
	}
	return r.converter.ObjectTypesSqlcToEntity(objectTypes), nil
}

func (r *Map) GetObjectsResponse(ctx context.Context, buildingID uuid.UUID) (entities.GetObjectsResponse, error) {
	buildingRow, err := r.q.GetBuildingByID(ctx, buildingID)
	if err != nil {
		return entities.GetObjectsResponse{}, fmt.Errorf("failed to get building: %w", err)
	}
	buildingEntity := entities.Building{
		ID:      buildingRow.ID,
		Name:    buildingRow.Name,
		Address: buildingRow.Address,
	}

	// Получаем этажи здания
	floorRows, err := r.q.GetFloorsByBuilding(ctx, buildingID)
	if err != nil {
		return entities.GetObjectsResponse{}, fmt.Errorf("failed to get floors: %w", err)
	}

	// Получаем объекты по зданию
	rowObjects, err := r.q.GetObjectsByBuilding(ctx, buildingID)
	if err != nil {
		return entities.GetObjectsResponse{}, fmt.Errorf("failed to get objects by building: %w", err)
	}

	// Собираем ID объектов для запроса дверей
	objectIDs := make([]uuid.UUID, len(rowObjects))
	for i, obj := range rowObjects {
		objectIDs[i] = obj.ID
	}

	rowDoors, err := r.q.GetDoorsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return entities.GetObjectsResponse{}, fmt.Errorf("failed to get doors: %w", err)
	}
	doorsMap := r.converter.DoorsSqlcToEntityMap(rowDoors)
	objects := r.converter.ObjectsSqlcToEntityByBuilding(rowObjects, doorsMap)

	// Группируем объекты по этажам
	floorObjectsMap := make(map[uuid.UUID][]entities.Object)
	for _, obj := range objects {
		floorObjectsMap[obj.Floor.ID] = append(floorObjectsMap[obj.Floor.ID], obj)
	}

	// Для каждого этажа получаем фон (background) и собираем данные
	var floorsWithData []entities.FloorWithData
	for _, fl := range floorRows {
		// Получаем фон по данному этажу
		bgRows, err := r.q.GetFloorBackground(ctx, fl.ID)
		if err != nil {
			return entities.GetObjectsResponse{}, fmt.Errorf(
				"failed to get floor background for floor %s: %w",
				fl.ID,
				err,
			)
		}
		background := r.converter.FloorBackgroundSqlcToEntityMany(bgRows)

		// Преобразуем строку этажа в сущность
		floorEntity := r.converter.FloorSqlcToEntity(fl)
		floorData := entities.FloorWithData{
			Floor:      floorEntity,
			Objects:    floorObjectsMap[fl.ID],
			Background: background,
		}
		floorsWithData = append(floorsWithData, floorData)
	}

	response := entities.GetObjectsResponse{
		Building: buildingEntity,
		Floors:   floorsWithData,
	}

	return response, nil
}

func (r *Map) GetObjectsByBuilding(ctx context.Context, buildingID uuid.UUID) ([]entities.Object, error) {
	rowObjects, err := r.q.GetObjectsByBuilding(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects by building: %w", err)
	}

	objectIDs := make([]uuid.UUID, len(rowObjects))
	for i, obj := range rowObjects {
		objectIDs[i] = obj.ID
	}

	rowDoors, err := r.q.GetDoorsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get doors: %w", err)
	}
	doorsMap := r.converter.DoorsSqlcToEntityMap(rowDoors)
	objects := r.converter.ObjectsSqlcToEntityByBuilding(rowObjects, doorsMap)
	return objects, nil
}

func (r *Map) UpdateObject(ctx context.Context, input entities.UpdateObjectInput) (entities.Object, error) {
	objectType, err := r.q.GetObjectTypeByName(ctx, string(input.ObjectType))
	if err != nil {
		return entities.Object{}, fmt.Errorf("failed to get object type: %w", err)
	}

	params := sqlc.UpdateObjectParams{
		Name:         input.Name,
		Alias:        input.Alias,
		Description:  sql.NullString{String: input.Description, Valid: input.Description != ""},
		ObjectTypeID: objectType.ID,
		ID:           input.ID,
	}
	rowObject, err := r.q.UpdateObject(ctx, params)
	if err != nil {
		return entities.Object{}, fmt.Errorf("failed to update object: %w", err)
	}

	rowDoors, err := r.q.GetDoorsByObjectIDs(ctx, []uuid.UUID{rowObject.ID})
	if err != nil {
		return entities.Object{}, fmt.Errorf("failed to get doors: %w", err)
	}
	doorsMap := r.converter.DoorsSqlcToEntityMap(rowDoors)

	description := ""
	if rowObject.Description.Valid {
		description = rowObject.Description.String
	}

	updatedObject := entities.Object{
		ID:          rowObject.ID,
		Name:        rowObject.Name,
		Alias:       rowObject.Alias,
		Description: description,
		X:           rowObject.X,
		Y:           rowObject.Y,
		Width:       rowObject.Width,
		Height:      rowObject.Height,
		ObjectType:  entities.ObjectType(objectType.Name),
		Doors:       doorsMap[rowObject.ID],
	}
	return updatedObject, nil
}
