package entities

import "github.com/google/uuid"

type UserContext struct {
	ID       uuid.UUID `json:"id"`
	Time     string    `json:"time"`
	Location *Location `json:"location"`
}

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
