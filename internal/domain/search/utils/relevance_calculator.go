package utils

import (
	"math"
	"strings"

	mapentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/entities"
	searchentities "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

const (
	baseMatchWeight  = 0.7
	distanceWeight   = 0.2
	popularityWeight = 0.1
)

type RelevanceCalculator struct {
	baseMatchWeight  float64
	distanceWeight   float64
	popularityWeight float64
}

func NewRelevanceCalculator() *RelevanceCalculator {
	return &RelevanceCalculator{
		baseMatchWeight:  baseMatchWeight,
		distanceWeight:   distanceWeight,
		popularityWeight: popularityWeight,
	}
}

func (c *RelevanceCalculator) Calculate(
	query string,
	obj mapentities.Object,
	ctx *searchentities.UserContext,
) float64 {
	baseScore := c.calculateBaseScore(query, obj)
	distanceScore := c.calculateDistanceScore(ctx, obj)
	return baseScore*c.baseMatchWeight + distanceScore*c.distanceWeight
}

func (c *RelevanceCalculator) calculateBaseScore(
	query string,
	obj mapentities.Object,
) float64 {
	score := 0.0

	if strings.Contains(strings.ToLower(string(obj.ObjectType)), strings.ToLower(query)) {
		score += 0.5
	}

	if strings.Contains(strings.ToLower(obj.Floor.Name), strings.ToLower(query)) {
		score += 0.3
	}

	if strings.Contains(strings.ToLower(obj.Alias), strings.ToLower(query)) {
		score += 0.2
	}

	return math.Min(score, 1.0)
}

func (c *RelevanceCalculator) calculateDistanceScore(
	ctx *searchentities.UserContext,
	obj mapentities.Object,
) float64 {
	if ctx == nil || ctx.Location == nil {
		return 0.0
	}

	dx := ctx.Location.X - obj.X
	dy := ctx.Location.Y - obj.Y
	distance := math.Sqrt(dx*dx + dy*dy)

	const maxDistance = 100.0
	normalizedDistance := math.Min(distance/maxDistance, 1.0)

	return 1 - normalizedDistance
}
