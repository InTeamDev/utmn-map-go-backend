package service

import (
	"context"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
	"github.com/google/uuid"
)

type RouteService interface {
	// GetRoute строит маршрут между точками
	// (первая точка - начальная, промежуточные, последняя - конечная).
	// Точки - ID Объектов.
	BuildRoute(ctx context.Context, start uuid.UUID, end uuid.UUID, waypoints []uuid.UUID) ([]entities.Edge, error)
	// Admin. SetPoint добавляет новую точку в граф.
	SetPoint(ctx context.Context, x, y float64) (uuid.UUID, error)
	// Admin. DeletePoint удаляет точку из графа.
	DeletePoint(ctx context.Context, id uuid.UUID) error
}
