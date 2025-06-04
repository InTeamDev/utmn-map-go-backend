package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/google/uuid"
)

type RouteRepository interface {
	CreateConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (entities.Edge, error)
	CreateIntersection(ctx context.Context, x, y float64, floorID uuid.UUID) (entities.Node, error)
	GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error)
	DeleteIntersection(ctx context.Context, buildingID uuid.UUID, id uuid.UUID) error
	GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]entities.Intersection, error)
}

type RouteService struct {
	repo RouteRepository
}

func NewRoute(repo RouteRepository) *RouteService {
	return &RouteService{
		repo: repo,
	}
}

func (r *RouteService) AddIntersection(ctx context.Context, x, y float64, floorID uuid.UUID) (uuid.UUID, error) {
	node, err := r.repo.CreateIntersection(ctx, x, y, floorID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create intersection in (%f;%f): %w", x, y, err)
	}
	return node.ID, nil
}

func (r *RouteService) GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]entities.Intersection, error) {
	intersections, err := r.repo.GetIntersections(ctx, buildingID)
	if err != nil {
		return []entities.Intersection{}, fmt.Errorf("get intersection: %w", err)
	}
	return intersections, nil
}

func (r *RouteService) AddConnection(
	ctx context.Context,
	fromID, toID uuid.UUID,
	weight float64,
) (entities.Edge, error) {
	conn, err := r.repo.CreateConnection(ctx, fromID, toID, weight)
	if err != nil {
		return entities.Edge{}, fmt.Errorf("create connection %w", err)
	}
	return entities.Edge{FromID: conn.FromID, ToID: conn.ToID, Weight: conn.Weight}, nil
}

func (r *RouteService) GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error) {
	connections, err := r.repo.GetConnections(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get connections for building %s: %w", buildingID, err)
	}
	return connections, nil
}

func (r *RouteService) DeleteIntersection(ctx context.Context, buildingID uuid.UUID, id uuid.UUID) error {
	err := r.repo.DeleteIntersection(ctx, buildingID, id)
	if err != nil {
		if err.Error() == "intersection not found" {
			return errors.New("intersection not found")
		}
		return fmt.Errorf("delete intersection %s: %w", id, err)
	}
	return nil
}
