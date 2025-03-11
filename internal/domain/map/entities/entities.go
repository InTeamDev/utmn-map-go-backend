package entities

import "github.com/google/uuid"

type GetObjectsRequest struct {
	BuildID uuid.UUID
	FloorID uuid.UUID
}
