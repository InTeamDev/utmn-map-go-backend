package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
)

type RouteRepository interface {
	CreateConnection(ctx context.Context, fromID, toID uuid.UUID, weight float64) (entities.Connection, error)
	CreateIntersection(ctx context.Context, req entities.AddIntersectionRequest) (entities.Node, error)
	GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error)
	DeleteIntersection(ctx context.Context, buildingID, id uuid.UUID) error
	GetIntersections(ctx context.Context, buildingID uuid.UUID) ([]entities.Node, error)
	GetDoors(ctx context.Context, buildingID uuid.UUID) ([]entities.Node, error)
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

func (r *RouteService) GetIntersections(ctx context.Context, buildID uuid.UUID) ([]entities.Node, error) {
	intersections, err := r.repo.GetIntersections(ctx, buildID)
	if err != nil {
		return []entities.Node{}, fmt.Errorf("get intersection: %w", err)
	}
	return intersections, nil
}

func (r *RouteService) AddConnection(
	ctx context.Context,
	fromID, toID uuid.UUID,
	weight float64,
) (entities.Connection, error) {
	conn, err := r.repo.CreateConnection(ctx, fromID, toID, weight)
	if err != nil {
		return entities.Connection{}, fmt.Errorf("create connection from %s to %s: %w", fromID, toID, err)
	}
	return entities.Connection{FromID: conn.FromID, ToID: conn.ToID, Weight: conn.Weight}, nil
}

func (r *RouteService) GetConnections(ctx context.Context, buildingID uuid.UUID) ([]entities.Connection, error) {
	connections, err := r.repo.GetConnections(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get connections for building %s: %w", buildingID, err)
	}
	return connections, nil
}

func (r *RouteService) GetDoors(ctx context.Context, buildingID uuid.UUID) ([]entities.Node, error) {
	doors, err := r.repo.GetDoors(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get doors for building %s: %w", buildingID, err)
	}
	return doors, nil
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

// BuildRoute теперь получает:
//
//	start, end – узлы (любые: дверь или пересечение),
//	waypoints – дополнительный список UUID-узлов, через которые нужно обязательно пройти.
//
// Алгоритм:
//  1. Сначала подгружаем из репозитория ВСЕ узлы (пересечения + двери) этого buildingID,
//     и формируем nodesMap для валидации наличия start/end/waypoints.
//  2. Проверяем, что start, end и все waypoints действительно принадлежат своему buildingID.
//  3. Получаем из репозитория ВСЕ связи (connections) этого buildingID.
//  4. Строим «маршрут» последовательно по отрезкам:
//     start → waypoints[0], waypoints[0] → waypoints[1], …, lastWaypoint → end.
//     Для каждого отрезка вызываем вспомогательную функцию buildShortestPath, которая
//     строит путь по алгоритму Дейкстры лишь по списку связей.
//  5. Склейка всех отрезков в единый []entities.Connection и возвращаем пользователю.
func (r *RouteService) BuildRoute(
	ctx context.Context,
	buildingID uuid.UUID,
	input entities.BuildRouteRequest,
) ([]entities.Connection, error) {
	// --- 2. Подгружаем все узлы (пересечения + двери) одного и того же buildingID.
	intersections, err := r.repo.GetIntersections(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("GetIntersections: %w", err)
	}
	doors, err := r.repo.GetDoors(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("GetDoors: %w", err)
	}
	// Заполняем nodesMap для быстрой проверки существования любого UUID:
	nodesMap := make(map[uuid.UUID]struct{}, len(intersections)+len(doors))
	for _, n := range intersections {
		nodesMap[n.ID] = struct{}{}
	}
	for _, n := range doors {
		nodesMap[n.ID] = struct{}{}
	}
	// Проверяем, что start, end и все waypoints есть в nodesMap:
	if _, ok := nodesMap[input.StartNodeID]; !ok {
		return nil, errors.New("start node не найден среди пересечений/дверей")
	}
	if _, ok := nodesMap[input.EndNodeID]; !ok {
		return nil, errors.New("end node не найден среди пересечений/дверей")
	}
	for _, wp := range input.Waypoints {
		if _, ok := nodesMap[wp]; !ok {
			return nil, fmt.Errorf("waypoint %s не найден среди пересечений/дверей", wp)
		}
	}

	// --- 3. Ловим из репозитория все связи (connections) для этого buildingID
	allConns, err := r.repo.GetConnections(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("get connections: %w", err)
	}

	// --- 4. Собираем порядок «посещений»: start → waypoints... → end
	order := make([]uuid.UUID, 0, len(input.Waypoints)+2)
	order = append(order, input.StartNodeID)
	order = append(order, input.Waypoints...)
	order = append(order, input.EndNodeID)

	var result []entities.Connection
	// Пробегаем по всем соседним парам в слайсе order:
	for i := 0; i < len(order)-1; i++ {
		segmentStart := order[i]
		segmentEnd := order[i+1]

		// Для каждого отрезка вызываем buildShortestPath, передавая полный список связей:
		edges, err := buildShortestPath(allConns, segmentStart, segmentEnd)
		if err != nil {
			return nil, fmt.Errorf("no route from %s to %s: %w", segmentStart, segmentEnd, err)
		}
		// Склеиваем все рёбра (edges) текущего отрезка в общий результат
		result = append(result, edges...)
	}

	return result, nil
}

// --- Вспомогательная функция: стандартная Дейкстра по списку связей.
// Возвращает список ребер (entities.Connection) в порядке следования кратчайшего пути.
func buildShortestPath(
	conns []entities.Connection,
	start, end uuid.UUID,
) ([]entities.Connection, error) {
	type neighbor struct {
		to     uuid.UUID
		weight float64
	}

	adj := make(map[uuid.UUID][]neighbor, len(conns))
	nodesSet := make(map[uuid.UUID]struct{}, len(conns)*2)
	for _, c := range conns {
		adj[c.FromID] = append(adj[c.FromID], neighbor{to: c.ToID, weight: c.Weight})
		adj[c.ToID] = append(adj[c.ToID], neighbor{to: c.FromID, weight: c.Weight})
		nodesSet[c.FromID] = struct{}{}
		nodesSet[c.ToID] = struct{}{}
	}

	// Карты для алгоритма:
	dist := make(map[uuid.UUID]float64, len(nodesSet))
	prev := make(map[uuid.UUID]uuid.UUID, len(nodesSet))
	unvisited := make(map[uuid.UUID]struct{}, len(nodesSet))
	const inf = 1e18

	for node := range nodesSet {
		dist[node] = inf
		unvisited[node] = struct{}{}
	}
	if _, ok := unvisited[start]; !ok {
		return nil, errors.New("start node not found in graph")
	}
	if _, ok := unvisited[end]; !ok {
		return nil, errors.New("end node not found in graph")
	}
	dist[start] = 0

	// Основной цикл Дейкстры (O(N^2) версия без heap):
	for len(unvisited) > 0 {
		// Находим вершину u из unvisited с минимальным dist[u]
		var u uuid.UUID
		minDist := inf
		for n := range unvisited {
			if d := dist[n]; d < minDist {
				minDist = d
				u = n
			}
		}

		// Если мы дошли до целевой вершины или minDist == inf — выходим
		if u == end || minDist == inf {
			break
		}
		delete(unvisited, u)

		// Расслабляем все рёбра из u
		for _, nb := range adj[u] {
			if _, seen := unvisited[nb.to]; !seen {
				continue
			}
			alt := dist[u] + nb.weight
			if alt < dist[nb.to] {
				dist[nb.to] = alt
				prev[nb.to] = u
			}
		}
	}

	// Если до конца не дошли — возвращаем ошибку
	if _, reached := prev[end]; !reached && start != end {
		return nil, errors.New("route not found")
	}

	// Восстанавливаем путь (список UUID) с конца в начало
	var path []uuid.UUID
	for cur := end; ; {
		path = append([]uuid.UUID{cur}, path...)
		if cur == start {
			break
		}
		p, ok := prev[cur]
		if !ok {
			return nil, errors.New("route reconstruction failed")
		}
		cur = p
	}

	// Конвертируем последовательность UUID в список ребер (entities.Connection)
	var edges []entities.Connection
	for i := 0; i < len(path)-1; i++ {
		u := path[i]
		v := path[i+1]
		// ищем вес ребра u → v
		w := 0.0
		for _, nb := range adj[u] {
			if nb.to == v {
				w = nb.weight
				break
			}
		}
		edges = append(edges, entities.Connection{
			FromID: u,
			ToID:   v,
			Weight: w,
		})
	}
	return edges, nil
}
