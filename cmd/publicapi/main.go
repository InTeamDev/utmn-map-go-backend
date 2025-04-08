package main

import (
	"flag"

	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/publicapi/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/publicapi.yaml", "Path to config file")
	flag.Parse()

	app.Run(configPath)
}
