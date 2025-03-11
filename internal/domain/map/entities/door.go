package entities

import (
	"github.com/google/uuid"
)

type Door struct {
	Id     uuid.UUID `json:"id"`
	X      float64   `json:"x"`
	Y      float64   `json:"y"`
	Width  float64   `json:"width"`
	Height float64   `json:"height"`
}
