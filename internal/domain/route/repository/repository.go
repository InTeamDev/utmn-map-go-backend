package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/repository/sqlc"
	"github.com/google/uuid"
)

type RouteConverter interface {
	ConnectionSqlcToEntity(i sqlc.Connection) entities.Connection
}

type RouteRepository struct {
	q         *sqlc.Queries
	converter *RouteConverterImpl
}

func NewRoute(db *sql.DB, converter RouteConverter) *RouteRepository {
	return &RouteRepository{
		q:         sqlc.New(db),
		converter: NewRouteConverter(),
	}
}

func (r *RouteRepository) CreateConnection(
	ctx context.Context,
	fromID, toID uuid.UUID,
	weight float64,
) (entities.Edge, error) {
	connection, err := r.q.CreateConnection(ctx, sqlc.CreateConnectionParams{
		FromID: fromID,
		ToID:   toID,
		Weight: weight,
	})
	if err != nil {
		return entities.Edge{}, err
	}
	edge := entities.Edge{
		FromID: connection.FromID,
		ToID:   connection.ToID,
		Weight: connection.Weight,
	}
	return edge, nil
}

func (r *RouteRepository) CreateIntersection(
	ctx context.Context,
	x, y float64,
	floorID uuid.UUID,
) (entities.Node, error) {
	intersection, err := r.q.CreateIntersection(ctx, sqlc.CreateIntersectionParams{
		ID:      uuid.New(),
		X:       x,
		Y:       y,
		FloorID: floorID,
	})
	if err != nil {
		return entities.Node{}, err
	}

	node := entities.Node{
		ID:   intersection.ID,
		X:    intersection.X,
		Y:    intersection.Y,
		Type: entities.NodeTypeIntersection,
	}
	return node, nil
}

func (r *RouteRepository) GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error) {
	sqlcConnections, err := r.q.GetConnections(ctx, buildingID)
	if err != nil {
		return nil, err
	}

	return r.converter.ConnectionsSqlcToEntity(sqlcConnections), nil
}

func (r *RouteRepository) DeleteIntersection(ctx context.Context, buildingID uuid.UUID, id uuid.UUID) error {
	err := r.q.DeleteIntersectionConnections(ctx, id)
	if err != nil {
		return fmt.Errorf("delete related connections: %w", err)
	}

	params := sqlc.DeleteIntersectionParams{
		IntersectionID: id,
		BuildingID:     buildingID,
	}
	err = r.q.DeleteIntersection(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("intersection not found")
		}
		return fmt.Errorf("delete intersection: %w", err)
	}

	return nil
}
