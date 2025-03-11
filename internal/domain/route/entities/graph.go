package entities

import "github.com/google/uuid"

type Graph struct {
	Nodes []uuid.UUID `json:"nodes"` // Список узлов
	Edges []Edges     `json:"edges"` // Список рёбер
}

type Edges struct {
	ID     uuid.UUID `json:"id"`
	From   uuid.UUID `json:"from"`
	To     uuid.UUID `json:"to"`
	Weight float64   `json:"weight"`
}
