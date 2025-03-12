package entities

import "github.com/google/uuid"

type Node struct {
	ID uuid.UUID `json:"id"`
	X  float64   `json:"x"`
	Y  float64   `json:"y"`
}

type Edge struct {
	ID     uuid.UUID `json:"id"`
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Weight float64   `json:"weight"`
}

type Graph struct {
	Nodes      map[uuid.UUID]Node `json:"nodes"` // Список узлов
	Edges      []Edge             `json:"edges"` // Список рёбер
	BuildingID uuid.UUID          // ID здания
}
