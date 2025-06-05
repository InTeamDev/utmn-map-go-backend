package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/repository/sqlc"
)

type RouteConverter interface {
	ConnectionSqlcToEntity(i sqlc.Connection) entities.Connection
	IntersectionSqlcToEntity(intersection sqlc.Intersection) entities.Intersection
	IntersectionsSqlcToEntity(intersections []sqlc.GetIntersectionsRow) []entities.Intersection
}

type RouteRepository struct {
	q         *sqlc.Queries
	db        *sql.DB
	converter *RouteConverterImpl
}

func NewRoute(db *sql.DB, converter RouteConverter) *RouteRepository {
	return &RouteRepository{
		q:         sqlc.New(db),
		db:        db,
		converter: NewRouteConverter(),
	}
}

func (r *RouteRepository) CreateConnection(
	ctx context.Context,
	fromID, toID uuid.UUID,
	weight float64,
) (entities.Connection, error) {
	connection, err := r.q.CreateConnection(ctx, sqlc.CreateConnectionParams{
		FromID: fromID,
		ToID:   toID,
		Weight: weight,
	})
	if err != nil {
		slog.Error(
			"failed to create connection",
			"slog",
			slog.String("fromID", fromID.String()),
			slog.String("toID", toID.String()),
			slog.Float64("weight", weight),
			slog.Any("error", err),
		)
		return entities.Connection{}, err
	}
	return entities.Connection{
		FromID: connection.FromID,
		ToID:   connection.ToID,
		Weight: connection.Weight,
	}, nil
}

func (r *RouteRepository) CreateIntersection(
	ctx context.Context,
	req entities.AddIntersectionRequest,
) (entities.Node, error) {
	id := req.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	intersection, err := r.q.CreateIntersection(ctx, sqlc.CreateIntersectionParams{
		ID:      id,
		X:       req.X,
		Y:       req.Y,
		FloorID: req.FloorID,
	})
	if err != nil {
		slog.Error(
			"failed to create intersection",
			"slog",
			slog.String("id", id.String()),
			slog.Float64("x", req.X),
			slog.Float64("y", req.Y),
			slog.String("floorID", req.FloorID.String()),
			slog.Any("error", err),
		)
		return entities.Node{}, fmt.Errorf("create intersection at (%f, %f): %w", req.X, req.Y, err)
	}

	node := entities.Node{
		ID:      intersection.ID,
		X:       intersection.X,
		Y:       intersection.Y,
		Type:    entities.NodeTypeIntersection,
		FloorID: intersection.FloorID,
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

func (r *RouteRepository) DeleteIntersection(ctx context.Context, buildingID, id uuid.UUID) error {
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

func (r *RouteRepository) GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]entities.Node, error) {
	intersections, err := r.q.GetIntersections(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get intersections by building: %w", err)
	}

	nodes := make([]entities.Node, len(intersections))
	for i, intersection := range intersections {
		nodes[i] = entities.Node{
			ID:      intersection.ID,
			X:       intersection.X,
			Y:       intersection.Y,
			Type:    entities.NodeTypeIntersection,
			FloorID: intersection.FloorID,
		}
	}
	return nodes, nil
}

func (r *RouteRepository) GetDoors(ctx context.Context, buildingID uuid.UUID) ([]entities.Node, error) {
	doors, err := r.q.ListDoorsByBuilding(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get doors by building: %w", err)
	}

	nodes := make([]entities.Node, len(doors))
	for i, door := range doors {
		nodes[i] = entities.Node{
			ID:      door.ID,
			X:       door.X,
			Y:       door.Y,
			Type:    entities.NodeTypeDoor,
			FloorID: door.FloorID,
		}
	}

	return nodes, nil
}
