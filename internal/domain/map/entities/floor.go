package entities

import "github.com/google/uuid"

type Floor struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
}

type BackgroundPoint struct {
	Order int     `json:"order"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type FloorBackgroundElement struct {
	ID     uuid.UUID         `json:"id"`
	Label  string            `json:"label"`
	ZIndex int               `json:"z_index"`
	Points []BackgroundPoint `json:"points"`
}
