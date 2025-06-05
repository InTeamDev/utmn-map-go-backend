package entities

import "github.com/google/uuid"

type AddIntersectionRequest struct {
	ID      uuid.UUID `json:"id"`
	X       float64   `json:"x"`
	Y       float64   `json:"y"`
	FloorID uuid.UUID `json:"floor_id"`
}
