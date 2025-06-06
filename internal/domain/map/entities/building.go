package entities

import "github.com/google/uuid"

type Building struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

type CreateBuildingInput struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}

type UpdateBuildingInput struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}
