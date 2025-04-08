package repository

import (
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	"github.com/google/uuid"
)

type MapConverterImpl struct{}

var _ MapConverter = (*MapConverterImpl)(nil)

func NewMapConverter() *MapConverterImpl {
	return &MapConverterImpl{}
}

func (mc *MapConverterImpl) floorSqlcToEntity(floors sqlc.Floor) entities.Floor {
	return entities.Floor{
		ID:    floors.ID,
		Name:  floors.Name,
		Alias: floors.Alias,
	}
}

func (mc *MapConverterImpl) FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor {
	result := make([]entities.Floor, 0, len(floors))
	for _, floor := range floors {
		result = append(result, mc.floorSqlcToEntity(floor))
	}
	return result
}

func (mc *MapConverterImpl) buildingSqlcToEntity(building sqlc.Building) entities.Building {
	return entities.Building{
		ID:      building.ID,
		Name:    building.Name,
		Address: building.Address,
	}
}

func (mc *MapConverterImpl) BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building {
	result := make([]entities.Building, 0, len(buildings))
	for _, building := range buildings {
		result = append(result, mc.buildingSqlcToEntity(building))
	}
	return result
}

func (mc *MapConverterImpl) ObjectSqlcToEntity(
	object sqlc.GetObjectsByBuildingRow,
	doors []entities.Door,
) entities.Object {
	description := ""
	if object.Description.Valid {
		description = object.Description.String
	}
	return entities.Object{
		ID:          object.ID,
		Name:        object.Name,
		Alias:       object.Alias,
		Description: description,
		X:           object.X,
		Y:           object.Y,
		Width:       object.Width,
		Height:      object.Height,
		ObjectType:  entities.ObjectType(object.ObjectType),
		Doors:       doors,
		Floor:       entities.Floor{ID: object.FloorID, Name: object.FloorName},
	}
}

func (mc *MapConverterImpl) doorSqlcToEntity(door sqlc.GetDoorsByObjectIDsRow) entities.Door {
	return entities.Door{
		ID:     door.ID,
		X:      door.X,
		Y:      door.Y,
		Width:  door.Width,
		Height: door.Height,
	}
}

func (mc *MapConverterImpl) DoorsSqlcToEntityMap(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door {
	result := make(map[uuid.UUID][]entities.Door)
	for _, door := range doors {
		result[door.ObjectID] = append(result[door.ObjectID], mc.doorSqlcToEntity(door))
	}
	return result
}

func (mc *MapConverterImpl) ObjectsSqlcToEntityByBuilding(
	objects []sqlc.GetObjectsByBuildingRow,
	doors map[uuid.UUID][]entities.Door,
) []entities.Object {
	result := make([]entities.Object, 0, len(objects))
	for _, object := range objects {
		result = append(result, mc.ObjectSqlcToEntity(object, doors[object.ID]))
	}
	return result
}

func (mc *MapConverterImpl) ObjectTypesSqlcToEntity(objectTypes []sqlc.ObjectType) []entities.ObjectType {
	result := make([]entities.ObjectType, 0, len(objectTypes))
	for _, objectType := range objectTypes {
		result = append(result, entities.ObjectType(objectType.Name))
	}
	return result
}
