package repository

import (
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/repository/sqlc"
)

type RouteConverterImpl struct{}

var _ RouteConverter = (*RouteConverterImpl)(nil)

func NewRouteConverter() *RouteConverterImpl {
	return &RouteConverterImpl{}
}

func (rc *RouteConverterImpl) IntersectionSqlcToEntity(i sqlc.Intersection) entities.Intersection {
	return entities.Intersection{
		ID:      i.ID,
		X:       i.X,
		Y:       i.Y,
		FloorID: i.FloorID,
	}
}

func (rc *RouteConverterImpl) ConnectionSqlcToEntity(i sqlc.Connection) entities.Connection {
	return entities.Connection{
		FromID: i.FromID,
		ToID:   i.ToID,
		Weight: i.Weight,
	}
}

func (rc *RouteConverterImpl) IntersectionsSqlcToEntity(
	intersections []sqlc.GetIntersectionsRow,
) []entities.Intersection {
	result := make([]entities.Intersection, 0, len(intersections))
	for _, intersection := range intersections {
		intersection := sqlc.Intersection{
			ID:      intersection.ID,
			X:       intersection.X,
			Y:       intersection.Y,
			FloorID: intersection.FloorID,
		}

		result = append(result, rc.IntersectionSqlcToEntity(intersection))
	}
	return result
}

func (rc *RouteConverterImpl) ConnectionsSqlcToEntity(connections []sqlc.Connection) []entities.Connection {
	result := make([]entities.Connection, 0, len(connections))
	for _, connection := range connections {
		result = append(result, rc.ConnectionSqlcToEntity(connection))
	}
	return result
}
