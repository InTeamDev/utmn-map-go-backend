package main

import (
	"math/rand"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/publicapi"
)

func init() {
	// Seed random number generator for auth codes
	rand.Seed(time.Now().UnixNano())
}

func main() {
	publicapi.Run()
}
