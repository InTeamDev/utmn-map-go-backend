package entities

import (
	"github.com/google/uuid"

	routeentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/entities"
)

type SyncAllDataRequest struct {
	Data SyncAllData `json:"data"`
}
type SyncAllData struct {
	ObjectTypes []ObjectTypeInfo `json:"object_types"`
	Buildings   []SyncBuildings  `json:"buildings"`
}

type SyncBuildings struct {
	ID      uuid.UUID    `json:"id"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Floors  []SyncFloors `json:"floors"`
}

type SyncFloors struct {
	ID            uuid.UUID                    `json:"id"`
	Name          string                       `json:"name"`
	Alias         string                       `json:"alias"`
	BuildingID    uuid.UUID                    `json:"building_id"`
	Objects       []Object                     `json:"objects"`
	Doors         []Door                       `json:"doors"`
	FloorPolygons []Polygon                    `json:"polygons"`
	Intersections []routeentities.Intersection `json:"intersections"`
	Connections   []routeentities.Connection   `json:"connections"`
}
