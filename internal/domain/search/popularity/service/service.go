package service

import "github.com/google/uuid"

type PopularityRanker struct {
	viewCounts  map[uuid.UUID]int
	clickCounts map[uuid.UUID]int
}

func NewPopularityRanker() *PopularityRanker {
	return &PopularityRanker{
		viewCounts:  make(map[uuid.UUID]int),
		clickCounts: make(map[uuid.UUID]int),
	}
}

func (r *PopularityRanker) UpdateStats(objectID uuid.UUID, viewed bool, clicked bool) {
	if viewed {
		r.viewCounts[objectID]++
	}
	if clicked {
		r.clickCounts[objectID]++
	}
}

func (r *PopularityRanker) GetPopularityScore(objectID uuid.UUID) float64 {
	views := r.viewCounts[objectID]
	clicks := r.clickCounts[objectID]

	double := 2
	if views == 0 {
		return float64(clicks * double)
	}
	return float64(clicks*2+views) / float64(views+1)
}
