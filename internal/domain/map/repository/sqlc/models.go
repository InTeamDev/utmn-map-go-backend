// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"github.com/google/uuid"
)

type Building struct {
	ID      uuid.UUID
	Name    string
	Address string
}

type Floor struct {
	ID         uuid.UUID
	Name       string
	Alias      string
	BuildingID uuid.UUID
}
