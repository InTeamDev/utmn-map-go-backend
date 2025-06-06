package entities

import "github.com/google/uuid"

type PolygonPoint struct {
	PolygonID uuid.UUID `json:"polygon_id"`
	Order     int32     `json:"order"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
}

type PolygonPointRequest struct {
	PointOrder int32   `json:"point_order"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
}


type DeletePolygonPointsRequest struct {
	PolygonID   uuid.UUID   `json:"polygon_id"`
	PointOrders []int32     `json:"point_orders"`
}
