package utils

import (
	"github.com/google/uuid"
)

func GenerateCacheKey(query string, userFloor uuid.UUID) string {
	cacheKey := query + ":" + userFloor.String()
	return cacheKey
}
