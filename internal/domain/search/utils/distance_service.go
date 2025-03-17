package utils

import (
	"math"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

type DistanceService struct {
	maxBuildingDistance float64
}

func NewDistanceService() *DistanceService {
	return &DistanceService{maxBuildingDistance: 100.0}
}

func (s *DistanceService) Calculate(
	userLoc searchentities.Location,
	obj mapentities.Object,
) float64 {
	dx := userLoc.X - obj.X
	dy := userLoc.Y - obj.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	return math.Max(0.1, 1-distance/s.maxBuildingDistance)
}
