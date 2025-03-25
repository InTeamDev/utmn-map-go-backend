package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	popularityservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/popularity/service"
	searchservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/utils"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/publicapi/http/handler"
	"github.com/InTeamDev/utmn-map-go-backend/internal/infrastructure/cache"
	"github.com/gin-gonic/gin"
)

const (
	ttl               = 5 * time.Minute
	readHeaderTimeout = 5 * time.Second
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP запросов",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время выполнения HTTP запроса",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Host == "" {
		return nil, errors.New("server.host is required")
	}
	if cfg.Server.Port == 0 {
		return nil, errors.New("server.port is required")
	}
	if cfg.Database.DSN == "" {
		return nil, errors.New("database.dsn is required")
	}

	return &cfg, nil
}

func run() int {
	var configPath string
	flag.StringVar(&configPath, "config", "config/publicapi.yaml", "Path to config file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting publicapi service...")

	if err := RunApp(ctx, configPath); err != nil {
		log.Printf("Service failed: %v", err)
		return 1
	}

	log.Println("Service stopped gracefully")
	return 0
}

func RunApp(ctx context.Context, configPath string) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initDatabase(config.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	mapService, mapRepository := initServices(db)
	searchService := initSearchService(mapRepository)

	router := gin.Default()

	// 🔥 Вставляем метрики
	router.Use(MetricsMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	publicAPI := handler.NewPublicAPI(mapService, searchService)
	publicAPI.RegisterRoutes(router)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return startServer(ctx, server)
}

func initDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func initServices(db *sql.DB) (*service.Map, *repository.Map) {
	mapConverter := repository.NewMapConverter()
	mapRepository := repository.NewMap(db, mapConverter)
	return service.NewMap(mapRepository), mapRepository
}

func initSearchService(mapRepository *repository.Map) *searchservice.SearchService {
	cache := cache.NewInMemorySearchCache(ttl)
	queryProcessor := utils.NewQueryProcessor("data/synonyms.json")
	popularityRanker := popularityservice.NewPopularityRanker()

	return searchservice.NewSearchService(
		cache,
		mapRepository,
		queryProcessor,
		popularityRanker,
	)
}

func startServer(ctx context.Context, server *http.Server) error {
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Starting server at %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal (context canceled)")
	case sig := <-stop:
		log.Printf("Received shutdown signal: %s", sig)
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), readHeaderTimeout)
	defer cancel()

	log.Println("Shutting down server...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server gracefully stopped")
	return nil
}

func main() {
	os.Exit(run())
}
