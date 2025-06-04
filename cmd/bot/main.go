package main

import (
	"flag"

	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/bot/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/bot.yaml", "Path to config file")
	flag.Parse()

	app.Run(configPath)
}
