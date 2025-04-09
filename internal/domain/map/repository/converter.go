package repository

import (
	"encoding/json"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository/sqlc"
	"github.com/google/uuid"
)

type MapConverterImpl struct{}

var _ MapConverter = (*MapConverterImpl)(nil)

func NewMapConverter() *MapConverterImpl {
	return &MapConverterImpl{}
}

func (mc *MapConverterImpl) FloorSqlcToEntity(f sqlc.Floor) entities.Floor {
	return entities.Floor{
		ID:    f.ID,
		Name:  f.Name,
		Alias: f.Alias,
	}
}

func (mc *MapConverterImpl) FloorsSqlcToEntity(floors []sqlc.Floor) []entities.Floor {
	result := make([]entities.Floor, 0, len(floors))
	for _, f := range floors {
		result = append(result, mc.FloorSqlcToEntity(f))
	}
	return result
}

func (mc *MapConverterImpl) buildingSqlcToEntity(b sqlc.Building) entities.Building {
	return entities.Building{
		ID:      b.ID,
		Name:    b.Name,
		Address: b.Address,
	}
}

func (mc *MapConverterImpl) BuildingsSqlcToEntity(buildings []sqlc.Building) []entities.Building {
	result := make([]entities.Building, 0, len(buildings))
	for _, b := range buildings {
		result = append(result, mc.buildingSqlcToEntity(b))
	}
	return result
}

func (mc *MapConverterImpl) ObjectSqlcToEntity(
	object sqlc.GetObjectsByBuildingRow,
	doors []entities.Door,
) entities.Object {
	return entities.Object{
		ID:          object.ID,
		Name:        object.Name,
		Alias:       object.Alias,
		Description: object.Description.String,
		X:           object.X,
		Y:           object.Y,
		Width:       object.Width,
		Height:      object.Height,
		ObjectType:  entities.ObjectType(object.ObjectType),
		Doors:       doors,
		Floor: entities.Floor{
			ID:    object.FloorID,
			Name:  object.FloorName,
			Alias: object.FloorName},
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

// Функция для преобразования фоновых элементов этажа.
// Предполагается, что поле Points типа []byte содержит JSON-массив точек,
// который преобразуем в []entities.BackgroundPoint.
func (mc *MapConverterImpl) FloorBackgroundSqlcToEntityMany(
	rows []sqlc.GetFloorBackgroundRow,
) []entities.FloorBackgroundElement {
	result := make([]entities.FloorBackgroundElement, 0, len(rows))
	for _, row := range rows {
		var points []entities.BackgroundPoint
		if err := json.Unmarshal(row.Points, &points); err != nil {
			// При ошибке можно логировать и пропускать этот элемент
			continue
		}
		element := entities.FloorBackgroundElement{
			ID:     row.ID,
			Label:  row.Label.String,
			ZIndex: int(row.ZIndex.Int32),
			Points: points,
		}
		result = append(result, element)
	}
	return result
}
