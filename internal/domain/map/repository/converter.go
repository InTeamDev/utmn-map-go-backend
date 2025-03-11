package repository

import (
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	"github.com/google/uuid"
)

type MapConverter struct{}

func NewMapConverter() *MapConverter {
	return &MapConverter{}
}

func (mc *MapConverter) ObjectSqlcToEntity(
	object sqlc.GetObjectsByBuildingAndFloorRow,
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
	}
}

func (mc *MapConverter) ObjectsSqlcToEntity(
	objects []sqlc.GetObjectsByBuildingAndFloorRow,
	doors map[uuid.UUID][]entities.Door,
) []entities.Object {
	result := make([]entities.Object, 0, len(objects))
	for _, object := range objects {
		result = append(result, mc.ObjectSqlcToEntity(object, doors[object.ID]))
	}
	return result
}

func (mc *MapConverter) DoorSqlcToEntity(door sqlc.GetDoorsByObjectIDsRow) entities.Door {
	return entities.Door{
		ID:     door.ID,
		X:      door.X,
		Y:      door.Y,
		Width:  door.Width,
		Height: door.Height,
	}
}

func (mc *MapConverter) DoorsSqlcToEntity(doors []sqlc.GetDoorsByObjectIDsRow) map[uuid.UUID][]entities.Door {
	result := make(map[uuid.UUID][]entities.Door)
	for _, door := range doors {
		result[door.ObjectID] = append(result[door.ObjectID], mc.DoorSqlcToEntity(door))
	}
	return result
}
