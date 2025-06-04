package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/InTeamDev/utmn-map-go-backend/config"
	authrepo "github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
	maprepository "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository"
	mapservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service"
	routerepository "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/repository"
	routeservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/route/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/adminapi/http/handler"
	"github.com/InTeamDev/utmn-map-go-backend/internal/middleware"
	"github.com/InTeamDev/utmn-map-go-backend/internal/server"
	"github.com/InTeamDev/utmn-map-go-backend/pkg/database"
)

const (
	ttl               = 5 * time.Minute
	readHeaderTimeout = 5 * time.Second
	maxAge            = 24 * time.Hour
)

func Run(configPath string) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting publicapi service...")

	if err := runApp(ctx, configPath); err != nil {
		log.Printf("Service failed: %v", err)
		return 1
	}

	log.Println("Service stopped gracefully")
	return 0
}

func runApp(ctx context.Context, configPath string) error {
	cfg, err := config.LoadAdminAPI(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := database.Init(cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	mapConverter := maprepository.NewMapConverter()
	mapRepository := maprepository.NewMap(db, mapConverter)
	mapService := mapservice.NewMap(mapRepository)
	routeConverter := routerepository.NewRouteConverter()
	routeRepository := routerepository.NewRoute(db, routeConverter)
	routeService := routeservice.NewRoute(routeRepository)
	repo := authrepo.NewInMemory()

	metrics := middleware.NewMetrics()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	}))

	router.Use(metrics.Middleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	adminAPI := handler.NewAdminAPI(mapService, routeService)
	adminAPI.RegisterRoutes(
		router,
		middleware.JWTAuth(middleware.JWTAuthConfig{Secret: []byte(cfg.JWTSecret), Repo: repo}),
	)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return server.StartServer(ctx, srv)
}
