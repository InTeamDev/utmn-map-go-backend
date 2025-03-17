package entities

import "github.com/google/uuid"

type SearchResult struct {
	ID         uuid.UUID `json:"id"`
	Relevance  float64   `json:"relevance"`
	Popularity float64   `json:"popularity"`
	Floor      string    `json:"floor"`
	Type       string    `json:"type"`
	Detail     string    `json:"detail"`
	X          float64   `json:"x"`
	Y          float64   `json:"y"`
	Distance   float64   `json:"distance,omitempty"`
}
