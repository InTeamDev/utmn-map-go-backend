package repository

import (
	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
)

type SearchRepository interface {
	GetAllObjects() ([]mapentities.Object, error)
	GetObjectsByFloor(floor string) ([]mapentities.Object, error)
}
