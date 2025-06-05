package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
)

type RouteRepository interface {
	CreateConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (entities.Edge, error)
	CreateIntersection(ctx context.Context, req entities.AddIntersectionRequest) (entities.Node, error)
	GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error)
	DeleteIntersection(ctx context.Context, buildingID, id uuid.UUID) error
	GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]entities.Intersection, error)
	GetNodeBuilding(ctx context.Context, nodeID uuid.UUID) (uuid.UUID, error)
}

type RouteService struct {
	repo RouteRepository
}

func NewRoute(repo RouteRepository) *RouteService {
	return &RouteService{
		repo: repo,
	}
}

func (r *RouteService) AddIntersection(
	ctx context.Context,
	req entities.AddIntersectionRequest,
) (entities.Node, error) {
	node, err := r.repo.CreateIntersection(ctx, req)
	if err != nil {
		return entities.Node{}, fmt.Errorf("create intersection: %w", err)
	}
	return node, nil
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
		return entities.Edge{}, fmt.Errorf("create connection from %s to %s: %w", fromID, toID, err)
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

func (r *RouteService) DeleteIntersection(ctx context.Context, buildingID, id uuid.UUID) error {
	err := r.repo.DeleteIntersection(ctx, buildingID, id)
	if err != nil {
		if err.Error() == "intersection not found" {
			return errors.New("intersection not found")
		}
		return fmt.Errorf("delete intersection %s: %w", id, err)
	}
	return nil
}

func (r *RouteService) BuildRoute(ctx context.Context, start uuid.UUID, end uuid.UUID, waypoints []uuid.UUID) ([]entities.Edge, error) {
	startBuilding, err := r.repo.GetNodeBuilding(ctx, start)
	if err != nil {
		return nil, fmt.Errorf("get start node building: %w", err)
	}
	endBuilding, err := r.repo.GetNodeBuilding(ctx, end)
	if err != nil {
		return nil, fmt.Errorf("get end node building: %w", err)
	}
	if startBuilding != endBuilding {
		return nil, errors.New("nodes belong to different buildings")
	}
	for _, wp := range waypoints {
		b, err := r.repo.GetNodeBuilding(ctx, wp)
		if err != nil {
			return nil, fmt.Errorf("get waypoint node building: %w", err)
		}
		if b != startBuilding {
			return nil, errors.New("nodes belong to different buildings")
		}
	}

	connections, err := r.repo.GetConnections(ctx, startBuilding)
	if err != nil {
		return nil, fmt.Errorf("get connections: %w", err)
	}

	order := make([]uuid.UUID, 0, len(waypoints)+2)
	order = append(order, start)
	order = append(order, waypoints...)
	order = append(order, end)

	var result []entities.Edge
	for i := 0; i < len(order)-1; i++ {
		edges, err := buildShortestPath(connections, order[i], order[i+1])
		if err != nil {
			return nil, err
		}
		result = append(result, edges...)
	}

	return result, nil
}

func buildShortestPath(conns []entities.Connection, start, end uuid.UUID) ([]entities.Edge, error) {
	type neighbor struct {
		to     uuid.UUID
		weight float64
	}
	adj := make(map[uuid.UUID][]neighbor)
	nodesSet := make(map[uuid.UUID]struct{})
	for _, c := range conns {
		adj[c.FromID] = append(adj[c.FromID], neighbor{to: c.ToID, weight: c.Weight})
		nodesSet[c.FromID] = struct{}{}
		nodesSet[c.ToID] = struct{}{}
	}

	dist := make(map[uuid.UUID]float64)
	prev := make(map[uuid.UUID]uuid.UUID)
	unvisited := make(map[uuid.UUID]struct{})
	const inf = 1e18

	for n := range nodesSet {
		dist[n] = inf
		unvisited[n] = struct{}{}
	}
	if _, ok := unvisited[start]; !ok {
		return nil, errors.New("start node not found")
	}
	if _, ok := unvisited[end]; !ok {
		return nil, errors.New("end node not found")
	}
	dist[start] = 0

	for len(unvisited) > 0 {
		// find node with smallest distance
		var u uuid.UUID
		minDist := inf
		for n := range unvisited {
			if d := dist[n]; d < minDist {
				minDist = d
				u = n
			}
		}

		if u == end || minDist == inf {
			break
		}
		delete(unvisited, u)

		for _, nb := range adj[u] {
			alt := dist[u] + nb.weight
			if alt < dist[nb.to] {
				dist[nb.to] = alt
				prev[nb.to] = u
			}
		}
	}

	if _, ok := prev[end]; !ok && start != end {
		return nil, errors.New("route not found")
	}

	// reconstruct path
	var path []uuid.UUID
	for u := end; ; {
		path = append([]uuid.UUID{u}, path...)
		if u == start {
			break
		}
		p, ok := prev[u]
		if !ok {
			return nil, errors.New("route not found")
		}
		u = p
	}

	// convert to edges
	edges := make([]entities.Edge, 0, len(path)-1)
	for i := 0; i < len(path)-1; i++ {
		weight := 0.0
		for _, nb := range adj[path[i]] {
			if nb.to == path[i+1] {
				weight = nb.weight
				break
			}
		}
		edges = append(edges, entities.Edge{FromID: path[i], ToID: path[i+1], Weight: weight})
	}
	return edges, nil
}
