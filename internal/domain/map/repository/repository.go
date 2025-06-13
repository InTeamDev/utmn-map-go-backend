package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
)

//go:generate mockgen -destination=mocks/mock_map_converter.go -package=mocks -source=repository.go MapConverter
type MapConverter interface {
	ObjectsSqlcToEntityByBuilding(
		objects []sqlc.GetObjectsByBuildingRow,
		doors map[uuid.UUID][]entities.Door,
	) []entities.Object
	ObjectSqlcToEntity(object sqlc.GetObjectsByBuildingRow, doors []entities.Door) entities.Object
	GetDoorsSqlcToEntity(doors []sqlc.GetDoorsByBuildingRow) []entities.GetDoorsResponse
	DoorsSqlcToEntityMap(doors []sqlc.Door) map[uuid.UUID][]entities.Door
	FloorSqlcToEntity(f sqlc.Floor) entities.Floor
	FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor
	BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building
	ObjectTypeSqlcToEntity(objectType sqlc.ObjectType) entities.ObjectTypeInfo
	ObjectTypesSqlcToEntity(objectTypes []sqlc.ObjectType) []entities.ObjectTypeInfo
	// Новая функция для конвертации background этажа
	FloorBackgroundSqlcToEntityMany(rows []sqlc.GetFloorBackgroundRow) []entities.FloorBackgroundElement
	SlicePolygonPointSqlcToEntity(
		rows []sqlc.ListPolygonPointsByPolygonIDRow,
	) []entities.PolygonPoint
	SlicePolygonSqlcToEntity(rows []sqlc.FloorPolygon) []entities.Polygon
}

type Map struct {
	q         *sqlc.Queries
	db        *sql.DB
	converter MapConverter
}

func NewMap(db *sql.DB, converter MapConverter) *Map {
	return &Map{
		q:         sqlc.New(db),
		db:        db,
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

func (r *Map) GetDoors(ctx context.Context, buildingID uuid.UUID) ([]entities.GetDoorsResponse, error) {
	rowDoors, err := r.q.GetDoorsByBuilding(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doors: %w", err)
	}

	doors := r.converter.GetDoorsSqlcToEntity(rowDoors)

	return doors, nil
}

func (r *Map) GetObjectTypes(ctx context.Context) ([]entities.ObjectTypeInfo, error) {
	objectTypes, err := r.q.GetObjectTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object types: %w", err)
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

func (r *Map) GetObjectTypeByID(
	ctx context.Context,
	id int32,
) (entities.ObjectTypeInfo, error) {
	dbObjectType, err := r.q.GetObjectTypeByID(ctx, id)
	if err != nil {
		return entities.ObjectTypeInfo{}, fmt.Errorf("database query failed: %w", err)
	}

	return r.converter.ObjectTypeSqlcToEntity(dbObjectType), nil
}

func (r *Map) GetObjectByID(
	ctx context.Context,
	objectID uuid.UUID,
) (entities.Object, error) {
	// 1. Получаем объект из базы
	dbObject, err := r.q.GetObjectByID(ctx, objectID)
	if err != nil {
		return entities.Object{}, fmt.Errorf("get object: %w", err)
	}

	// 2. Получаем двери для этого объекта
	// Создаем слайс с одним ID
	objectIDs := []uuid.UUID{objectID}
	dbDoors, err := r.q.GetDoorsByObjectIDs(ctx, objectIDs)
	if err != nil {
		return entities.Object{}, fmt.Errorf("get doors: %w", err)
	}

	// 3. Конвертируем данные
	doorsMap := r.converter.DoorsSqlcToEntityMap(dbDoors)

	dbfloor, err := r.q.GetFloorByID(ctx, dbObject.FloorID)
	if err != nil {
		return entities.Object{}, fmt.Errorf("get floor: %w", err)
	}

	// Создаем временный объект нужного типа для конвертера
	tempObj := sqlc.GetObjectsByBuildingRow{
		ID:          dbObject.ID,
		Name:        dbObject.Name,
		Alias:       dbObject.Alias,
		Description: dbObject.Description,
		X:           dbObject.X,
		Y:           dbObject.Y,
		Width:       dbObject.Width,
		Height:      dbObject.Height,
		ObjectType:  dbObject.ObjectTypeID,
		FloorID:     dbfloor.ID,
		FloorName:   dbfloor.Name,
	}

	// 4. Вызываем конвертер
	return r.converter.ObjectSqlcToEntity(tempObj, doorsMap[objectID]), nil
}

func (r *Map) CreateObject(
	ctx context.Context,
	floorID uuid.UUID,
	input entities.CreateObjectInput,
) (entities.Object, error) {
	objectID := input.ID
	if objectID == uuid.Nil {
		objectID = uuid.New()
	}

	params := sqlc.CreateObjectParams{
		ID:           objectID,
		FloorID:      floorID,
		Name:         input.Name,
		Alias:        input.Alias,
		Description:  sql.NullString{String: input.Description, Valid: input.Description != ""},
		X:            input.X,
		Y:            input.Y,
		Width:        input.Width,
		Height:       input.Height,
		ObjectTypeID: input.ObjectTypeID,
	}

	rowObject, err := r.q.CreateObject(ctx, params)
	if err != nil {
		return entities.Object{}, fmt.Errorf("create object: %w", err)
	}

	description := ""
	if rowObject.Description.Valid {
		description = rowObject.Description.String
	}

	createdObject := entities.Object{
		ID:           rowObject.ID,
		Name:         rowObject.Name,
		Alias:        rowObject.Alias,
		Description:  description,
		X:            rowObject.X,
		Y:            rowObject.Y,
		Width:        rowObject.Width,
		Height:       rowObject.Height,
		ObjectTypeID: rowObject.ObjectTypeID,
		Doors:        nil,
	}

	return createdObject, nil
}

func (r *Map) UpdateObject(
	ctx context.Context,
	id uuid.UUID,
	input entities.UpdateObjectInput,
) (entities.Object, error) {
	params := sqlc.UpdateObjectParams{
		ID:           id,
		Name:         sqlNullString(input.Name),
		Alias:        sqlNullString(input.Alias),
		Description:  sqlNullString(input.Description),
		X:            sqlNullFloat64(input.X),
		Y:            sqlNullFloat64(input.Y),
		Width:        sqlNullFloat64(input.Width),
		Height:       sqlNullFloat64(input.Height),
		ObjectTypeID: sqlNullInt32(input.ObjectTypeID),
	}
	rowObject, err := r.q.UpdateObject(ctx, params)
	if err != nil {
		return entities.Object{}, fmt.Errorf("update object: %w", err)
	}

	description := ""
	if rowObject.Description.Valid {
		description = rowObject.Description.String
	}

	updatedObject := entities.Object{
		ID:           rowObject.ID,
		Name:         rowObject.Name,
		Alias:        rowObject.Alias,
		Description:  description,
		X:            rowObject.X,
		Y:            rowObject.Y,
		Width:        rowObject.Width,
		Height:       rowObject.Height,
		ObjectTypeID: rowObject.ObjectTypeID,
	}
	return updatedObject, nil
}

func (r *Map) DeleteObject(ctx context.Context, objectID uuid.UUID) error {
	return r.q.DeleteObject(ctx, objectID)
}

func (r *Map) GetDoor(
	ctx context.Context,
	buildingID uuid.UUID,
	doorID uuid.UUID,
) (entities.Door, error) {
	dbDoor, err := r.q.GetDoor(ctx, sqlc.GetDoorParams{
		Doorid:     doorID,
		Buildingid: buildingID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Door{}, sql.ErrNoRows
		}
		return entities.Door{}, fmt.Errorf("get door: %w", err)
	}

	result := entities.Door{
		ID:       dbDoor.ID,
		X:        dbDoor.X,
		Y:        dbDoor.Y,
		Width:    dbDoor.Width,
		Height:   dbDoor.Height,
		ObjectID: dbDoor.ObjectID,
	}

	return result, nil
}

func (r *Map) UpdateDoor(
	ctx context.Context,
	buildingID uuid.UUID,
	doorID uuid.UUID,
	input entities.Door,
) (entities.Door, error) {
	params := sqlc.UpdateDoorParams{
		DoorID:     doorID,
		BuildingID: buildingID,
		X:          input.X,
		Y:          input.Y,
		Width:      input.Width,
		Height:     input.Height,
		ObjectID:   input.ObjectID,
	}

	dbDoor, err := r.q.UpdateDoor(ctx, params)
	if err != nil {
		return entities.Door{}, err
	}

	return entities.Door{
		ID:       dbDoor.ID,
		X:        dbDoor.X,
		Y:        dbDoor.Y,
		Width:    dbDoor.Width,
		Height:   dbDoor.Height,
		ObjectID: dbDoor.ObjectID,
	}, nil
}

func (r *Map) CreateBuilding(ctx context.Context, input entities.CreateBuildingInput) (entities.Building, error) {
	id := input.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	params := sqlc.CreateBuildingParams{
		ID:      id,
		Name:    input.Name,
		Address: input.Address,
	}

	building, err := r.q.CreateBuilding(ctx, params)
	if err != nil {
		return entities.Building{}, fmt.Errorf("create building: %w", err)
	}

	result := entities.Building{
		ID:      building.ID,
		Name:    building.Name,
		Address: building.Address,
	}

	return result, nil
}

func (r *Map) DeleteBuilding(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteBuilding(ctx, id)
}

func (r *Map) UpdateBuilding(
	ctx context.Context,
	id uuid.UUID,
	input entities.UpdateBuildingInput,
) (entities.Building, error) {
	params := sqlc.UpdateBuildingParams{
		ID:      id,
		Name:    sql.NullString{String: input.Name, Valid: input.Name != ""},
		Address: sql.NullString{String: input.Address, Valid: input.Address != ""},
	}
	b, err := r.q.UpdateBuilding(ctx, params)
	if err != nil {
		return entities.Building{}, fmt.Errorf("update building: %w", err)
	}
	return entities.Building{
		ID:      b.ID,
		Name:    b.Name,
		Address: b.Address,
	}, nil
}

func (r *Map) GetBuildingByID(ctx context.Context, id uuid.UUID) (entities.Building, error) {
	b, err := r.q.GetBuildingByID(ctx, id)
	if err != nil {
		return entities.Building{}, fmt.Errorf("get building by id: %w", err)
	}
	return entities.Building{
		ID:      b.ID,
		Name:    b.Name,
		Address: b.Address,
	}, nil
}

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func sqlNullFloat64(i *float64) sql.NullFloat64 {
	if i == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *i, Valid: true}
}

func (r *Map) CreatePolygon(ctx context.Context, polygon entities.Polygon) (entities.Polygon, error) {
	id := polygon.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	row, err := r.q.CreatePolygon(ctx, sqlc.CreatePolygonParams{
		ID:      id,
		FloorID: polygon.FloorID,
		Label:   sql.NullString{String: polygon.Label, Valid: polygon.Label != ""},
		ZIndex:  sql.NullInt32{Int32: polygon.ZIndex, Valid: true},
	})
	if err != nil {
		return entities.Polygon{}, fmt.Errorf("create polygon: %w", err)
	}

	return entities.Polygon{
		ID:      row.ID,
		FloorID: row.FloorID,
		Label:   row.Label.String,
		ZIndex:  row.ZIndex.Int32,
	}, nil
}

func (r *Map) CreatePolygonPoint(
	ctx context.Context,
	polygonID uuid.UUID,
	order int32,
	x, y float64,
) (entities.PolygonPoint, error) {
	point, err := r.q.CreatePolygonPoint(ctx, sqlc.CreatePolygonPointParams{
		PolygonID:  polygonID,
		PointOrder: order,
		X:          x,
		Y:          y,
	})
	if err != nil {
		return entities.PolygonPoint{}, fmt.Errorf("create polygon point: %w", err)
	}
	return entities.PolygonPoint{
		PolygonID: point.PolygonID,
		Order:     point.PointOrder,
		X:         point.X,
		Y:         point.Y,
	}, nil
}

func (r *Map) CreateFloor(ctx context.Context, buildingID uuid.UUID, floor entities.Floor) error {
	id := floor.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	err := r.q.CreateFloor(ctx, sqlc.CreateFloorParams{
		ID:         id,
		Name:       floor.Name,
		Alias:      floor.Alias,
		BuildingID: buildingID,
	})
	if err != nil {
		return fmt.Errorf("create floor: %w", err)
	}
	return nil
}

func (r *Map) CreateDoor(ctx context.Context, objectID uuid.UUID, door entities.Door) (entities.Door, error) {
	id := door.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	params := sqlc.CreateDoorParams{
		ID:       id,
		X:        door.X,
		Y:        door.Y,
		Width:    door.Width,
		Height:   door.Height,
		ObjectID: door.ObjectID,
	}

	rowDoor, err := r.q.CreateDoor(ctx, params)
	if err != nil {
		return entities.Door{}, fmt.Errorf("create door: %w", err)
	}

	createdDoor := entities.Door{
		ID:       rowDoor.ID,
		X:        rowDoor.X,
		Y:        rowDoor.Y,
		Width:    rowDoor.Width,
		Height:   rowDoor.Height,
		ObjectID: rowDoor.ObjectID,
	}

	return createdDoor, nil
}

// GetDoorFloorPairs retrieves a map[doorID]floorID.
func (r *Map) GetDoorFloorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	pairs, err := r.q.GetDoorFloorPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get doors by building: %w", err)
	}
	doorsMap := make(map[uuid.UUID]uuid.UUID, len(pairs))
	for _, pair := range pairs {
		doorsMap[pair.DoorID] = pair.FloorID
	}
	return doorsMap, nil
}

// GetObjectDoorPairs retrieves a map[objectID]doorID.
func (r *Map) GetObjectDoorPairs(ctx context.Context) (map[uuid.UUID]uuid.UUID, error) {
	pairs, err := r.q.GetObjectDoorPairs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object door pairs: %w", err)
	}
	objectDoorMap := make(map[uuid.UUID]uuid.UUID, len(pairs))
	for _, pair := range pairs {
		objectDoorMap[pair.ObjectID] = pair.DoorID
	}
	return objectDoorMap, nil
}

func (r *Map) DeletePolygonPoints(ctx context.Context, request entities.DeletePolygonPointsRequest) error {
	return r.q.DeletePolygonPoints(ctx, sqlc.DeletePolygonPointsParams{
		PolygonID: request.PolygonID,
		Column2:   request.PointOrders,
	})
}

func (r *Map) GetPolygonByID(ctx context.Context, id uuid.UUID) (entities.Polygon, error) {
	dbPolygon, err := r.q.GetPolygonByID(ctx, id)
	if err != nil {
		return entities.Polygon{}, err
	}
	pts, err := r.q.ListPolygonPointsByPolygonID(ctx, id)
	if err != nil {
		return entities.Polygon{}, err
	}
	points := r.converter.SlicePolygonPointSqlcToEntity(pts)
	polygon := entities.Polygon{
		ID:      dbPolygon.ID,
		FloorID: dbPolygon.FloorID,
		Label:   dbPolygon.Label.String,
		ZIndex:  int32(dbPolygon.ZIndex.Int32),
		Points:  points,
	}
	return polygon, nil
}

func (r *Map) GetPolygonsByFloorID(ctx context.Context, floorID uuid.UUID) ([]entities.Polygon, error) {
	dbPolygons, err := r.q.GetPolygonsByFloorID(ctx, floorID)
	if err != nil {
		return nil, err
	}
	return r.converter.SlicePolygonSqlcToEntity(dbPolygons), nil
}

func (r *Map) ChangePolygonPoint(ctx context.Context, req entities.ChangePolygonPointRequest) error {
	return r.q.ChangePolygonPoint(ctx, sqlc.ChangePolygonPointParams{
		PolygonID:     req.PolygonID,
		OldPointOrder: req.OldPointOrder,
		PointOrder:    req.NewPointOrder,
		X:             req.X,
		Y:             req.Y,
	})
}
