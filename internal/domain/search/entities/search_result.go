package entities

import (
	"github.com/google/uuid"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type SearchRequest struct {
	Query      string                   `json:"query"`
	Limit      int                      `json:"limit"`
	Offset     int                      `json:"offset"`
	BuildingID uuid.UUID                `json:"building_id"`
	Categories []mapentities.ObjectType `json:"categories"`
}

type SearchResult struct {
	ObjectID     uuid.UUID `json:"object_id"`
	ObjectTypeID int       `json:"object_type_id"`
	Preview      string    `json:"preview"`
}
