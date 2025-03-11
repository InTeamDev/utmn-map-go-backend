package entities

import "github.com/google/uuid"

type Floor struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Alias   string    `json:"alias"`
	Objects []Object  `json:"objects"`
}
