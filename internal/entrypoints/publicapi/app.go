package publicapi

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

	"github.com/InTeamDev/utmn-map-go-backend/config"
	authService "github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	userRepository "github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/repository"
	userService "github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/publicapi/http/handler"
	"github.com/InTeamDev/utmn-map-go-backend/internal/tgbot"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Run starts the public API entrypoint
func Run() {
	// Load configuration
	cfg, err := config.New("config/publicapi.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Init repositories
	userRepo := userRepository.NewPostgresRepository(db)

	// Init services
	userSvc := userService.NewService(userRepo)
	authSvc := authService.NewAuthService(userSvc, cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Create and start Telegram bot
	botConfig := tgbot.Config{
		Token:          cfg.TGBot.Token,
		DevelopersChat: cfg.TGBot.DevelopersChatID,
	}
	telegramBot, err := tgbot.NewTelegramBot(botConfig, authSvc, userSvc)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}
	go telegramBot.Start()

	// Create HTTP server
	router := gin.Default()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	// Register handlers
	authHandler := handler.NewAuthHandler(authSvc, telegramBot)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Register API handlers
	authHandler.RegisterHandlers(router)

	// Start HTTP server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
