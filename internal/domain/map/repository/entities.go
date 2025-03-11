package repository

import "github.com/google/uuid"

type Object struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias"`
	Description string    `json:"description"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Width       float64   `json:"width"`
	Height      float64   `json:"height"`
	ObjectType  string    `json:"object_type"`
}
