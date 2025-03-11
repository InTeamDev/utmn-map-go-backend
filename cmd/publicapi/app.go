package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/repository"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/map/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/publicapi/http/handler"
	"github.com/gin-gonic/gin"
)

const readHeaderTimeout = 5 * time.Second

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

	mapService := initServices(db)

	router := gin.Default()
	publicAPI := handler.NewPublicAPI(mapService)
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

func initServices(db *sql.DB) *service.Map {
	mapConverter := repository.NewMapConverter()
	mapRepository := repository.NewMap(db, mapConverter)
	return service.NewMap(mapRepository)
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
