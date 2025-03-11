package entities

import "github.com/google/uuid"

type Building struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Floors  []Floor   `json:"floors"`
}
