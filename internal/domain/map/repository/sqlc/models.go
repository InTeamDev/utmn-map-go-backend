// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package sqlc

import (
	"database/sql"

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

type Object struct {
	ID           uuid.UUID
	Name         string
	Alias        string
	Description  sql.NullString
	X            float64
	Y            float64
	Width        float64
	Height       float64
	ObjectTypeID int32
	FloorID      uuid.UUID
}

type ObjectType struct {
	ID    int32
	Name  string
	Alias string
}
