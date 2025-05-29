package main

import (
	"flag"

	_ "github.com/InTeamDev/utmn-map-go-backend/cmd/adminapi/docs" // Импорт сгенерированной документации
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/adminapi/app"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/adminapi.yaml", "Path to config file")
	flag.Parse()

	app.Run(configPath)
}
