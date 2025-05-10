package entities

import "github.com/google/uuid"

// Edge - ребро графа.
type Edge struct {
	FromID uuid.UUID `json:"from_id"`
	ToID   uuid.UUID `json:"to_id"`
	Weight float64   `json:"weight"`
}

type NodeType string

const (
	NodeTypeIntersection NodeType = "intersection"
	NodeTypeDoor         NodeType = "door"
)

// Node - узел графа.
type Node struct {
	// Intersection ID || Door ID
	ID   uuid.UUID `json:"id"`
	X    float64   `json:"x"`
	Y    float64   `json:"y"`
	Type NodeType  `json:"type"`
}

type Graph struct {
	Nodes      map[uuid.UUID]Node `json:"nodes"` // Список узлов
	Edges      []Edge             `json:"edges"` // Список рёбер
	BuildingID uuid.UUID          `json:"building_id"`
}
