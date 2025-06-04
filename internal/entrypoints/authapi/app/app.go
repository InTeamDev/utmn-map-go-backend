package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/InTeamDev/utmn-map-go-backend/config"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
	authservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/authapi/http/handler"
	"github.com/InTeamDev/utmn-map-go-backend/internal/server"
)

const readHeaderTimeout = 5 * time.Second

type BotClient struct {
	url      string
	clientID string
	token    string
}

func (b *BotClient) SendMessage(chatID int64, msg string) error {
	payload := fmt.Sprintf(`{"telegram_user_id":%d,"message":"%s"}`, chatID, msg)
	req, err := http.NewRequest("POST", b.url+"/api/message", strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(b.clientID, b.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bot status %d", resp.StatusCode)
	}
	return nil
}

func Run(configPath string) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting authapi service...")

	if err := runApp(ctx, configPath); err != nil {
		log.Printf("Service failed: %v", err)
		return 1
	}
	log.Println("Service stopped gracefully")
	return 0
}

func runApp(ctx context.Context, configPath string) error {
	cfg, err := config.LoadAuthAPI(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	repo := repository.NewInMemory()
	bot := &BotClient{url: cfg.BotClient.URL, clientID: cfg.BotClient.ClientID, token: cfg.BotClient.AccessToken}
	svc := authservice.New(repo, bot, []byte(cfg.JWTSecret))

	r := gin.Default()
	api := handler.NewAuthAPI(svc, repo, []byte(cfg.JWTSecret), cfg.Auth.ClientID, cfg.Auth.AccessToken)
	api.RegisterRoutes(r)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	return server.StartServer(ctx, srv)
}
