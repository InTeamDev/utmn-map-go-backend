package entities

import "github.com/google/uuid"

type Connection struct {
	FromID uuid.UUID `json:"from_id"`
	ToID   uuid.UUID `json:"to_id"`
	Weight float64   `json:"weight"`
}

type Intersection struct {
	ID      uuid.UUID `json:"id"`
	X       float64   `json:"x"`
	Y       float64   `json:"y"`
	FloorID uuid.UUID `json:"floor_id"`
}

type NodeType string

const (
	NodeTypeIntersection NodeType = "intersection"
	NodeTypeDoor         NodeType = "door"
)

// Node - узел графа.
type Node struct {
	// Intersection ID || Door ID
	ID      uuid.UUID `json:"id"`
	X       float64   `json:"x"`
	Y       float64   `json:"y"`
	Type    NodeType  `json:"type"`
	FloorID uuid.UUID `json:"floor_id"`
}
