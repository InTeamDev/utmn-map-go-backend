package entities

import "github.com/google/uuid"

type Polygon struct {
	ID      uuid.UUID      `json:"id"`
	FloorID uuid.UUID      `json:"floor_id"`
	Label   string         `json:"label"`
	ZIndex  int32          `json:"z_index"`
	Points  []PolygonPoint `json:"points"`
}

type FloorPolygon struct {
	ID      uuid.UUID `json:"id"`
	FloorID uuid.UUID `json:"floor_id"`
	Label   string    `json:"label"`
	ZIndex  int       `json:"z_index"`
}
