package repository

import (
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	"github.com/google/uuid"
)

type GetObjectsRequest struct {
	BuildingID uuid.UUID
	FloorID    uuid.UUID
}

type GetObjectsResponse struct {
	Objects []entities.Object
}
