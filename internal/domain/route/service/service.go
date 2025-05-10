package service

import (
	"context"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/google/uuid"
)

type RouteRepository interface {
	CreateConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (entities.Edge, error)
	CreateIntersection(ctx context.Context, x, y float64) (entities.Node, error)
}

type RouteService struct {
	repo RouteRepository
}

func NewRoute(repo RouteRepository) *RouteService {
	return &RouteService{
		repo: repo,
	}
}

func (r *RouteService) AddIntersection(ctx context.Context, x, y float64) (uuid.UUID, error) {
	node, err := r.repo.CreateIntersection(ctx, x, y)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create intersection in (%f;%f): %w", x, y, err)
	}
	return node.ID, nil
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
