package entities

import (
	"github.com/google/uuid"
)

type Door struct {
	ID       uuid.UUID `json:"id"`
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	Width    float64   `json:"width"`
	Height   float64   `json:"height"`
	ObjectID uuid.UUID `json:"object_id"`
}

type UpdateDoorInput struct {
	X        *float64   `json:"x,omitempty"`
	Y        *float64   `json:"y,omitempty"`
	Width    *float64   `json:"width,omitempty"`
	Height   *float64   `json:"height,omitempty"`
	ObjectID *uuid.UUID `json:"object_id,omitempty"`
}
