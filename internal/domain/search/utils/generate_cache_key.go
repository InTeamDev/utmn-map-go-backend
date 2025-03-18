package utils

import (
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/entities"
	"github.com/google/uuid"
)

func GenerateCacheKey(query string, userFloor uuid.UUID, userContext *entities.UserContext) string {
	cacheKey := query + ":" + userFloor.String()
	if userContext != nil && userContext.Location != nil {
		cacheKey += fmt.Sprintf(":%f:%f", userContext.Location.X, userContext.Location.Y)
	}
	return cacheKey
}
