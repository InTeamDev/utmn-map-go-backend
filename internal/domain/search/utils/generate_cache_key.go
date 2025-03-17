package utils

import (
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
)

func GenerateCacheKey(query, userFloor string, userContext *entities.UserContext) string {
	cacheKey := query + ":" + userFloor
	if userContext != nil && userContext.Location != nil {
		cacheKey += fmt.Sprintf(":%f:%f", userContext.Location.X, userContext.Location.Y)
	}
	return cacheKey
}
