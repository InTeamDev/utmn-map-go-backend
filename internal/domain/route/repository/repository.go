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
) (entities.Edge, error) {
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
		return entities.Edge{}, err
	}
	return entities.Edge{
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

func (r *RouteRepository) GetIntersections(ctx context.Context, buildID uuid.UUID) ([]entities.Intersection, error) {
	intersections, err := r.q.GetIntersections(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get intersections by building: %w", err)
	}

	return r.converter.IntersectionsSqlcToEntity(intersections), nil
}

func (r *RouteRepository) GetNodeBuilding(ctx context.Context, nodeID uuid.UUID) (uuid.UUID, error) {
	const query = `SELECT f.building_id
                FROM intersections i
                JOIN floors f ON i.floor_id = f.id
                WHERE i.id = $1
                UNION
                SELECT f.building_id
                FROM doors d
                JOIN objects o ON d.object_id = o.id
                JOIN floors f ON o.floor_id = f.id
                WHERE d.id = $1
                LIMIT 1`

	var bID uuid.UUID
	if err := r.db.QueryRowContext(ctx, query, nodeID).Scan(&bID); err != nil {
		return uuid.Nil, err
	}
	return bID, nil
}
