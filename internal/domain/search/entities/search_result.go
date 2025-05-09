package entities

import (
	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/google/uuid"
)

type SearchRequest struct {
	Query      string                   `json:"query"`
	Limit      int                      `json:"limit"`
	Offset     int                      `json:"offset"`
	BuildingID uuid.UUID                `json:"building_id"`
	Categories []mapentities.ObjectType `json:"categories"`
}

type SearchResult struct {
	ObjectID uuid.UUID `json:"object_id"`
	Category string    `json:"category"`
	Preview  string    `json:"preview"`
}
