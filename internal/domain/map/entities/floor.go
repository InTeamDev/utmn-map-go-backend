package entities

import "github.com/google/uuid"

type Floor struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Alias string    `json:"alias"`
}
