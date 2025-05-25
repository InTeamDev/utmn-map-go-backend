package entities

import "github.com/google/uuid"

type PolygonPoint struct {
	ID        int32     `json:"id"`
	PolygonID uuid.UUID `json:"polygon_id"`
	Order     int32     `json:"point_order"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
}

type PolygonPointRequest struct {
	PointOrder int32   `json:"point_order" binding:"required"`
	X          float64 `json:"x"           binding:"required"`
	Y          float64 `json:"y"           binding:"required"`
}
