package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/telebot.v3"

	"github.com/InTeamDev/utmn-map-go-backend/config"
	"github.com/InTeamDev/utmn-map-go-backend/internal/entrypoints/bot/http/handler"
	"github.com/InTeamDev/utmn-map-go-backend/internal/server"
)

const readHeaderTimeout = 5 * time.Second

func Run(configPath string) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting bot service...")

	if err := runApp(ctx, configPath); err != nil {
		log.Printf("Service failed: %v", err)
		return 1
	}

	log.Println("Service stopped gracefully")
	return 0
}

func runApp(ctx context.Context, configPath string) error {
	cfg, err := config.LoadBot(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("init telebot: %w", err)
	}

	registerCommands(bot, cfg)

	go bot.Start()

	router := gin.Default()

	h := handler.NewBotHandler(bot)
	api := router.Group("/api", gin.BasicAuth(gin.Accounts{
		cfg.Auth.ClientID: cfg.Auth.AccessToken,
	}))
	api.POST("/message", h.SendMessage)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	return server.StartServer(ctx, srv)
}

func registerCommands(bot *telebot.Bot, cfg *config.Bot) {
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Привет! Для регистрации выполни команду /register")
	})

	bot.Handle("/register", func(c telebot.Context) error {
		if cfg.Backend.URL == "" {
			return c.Send("Backend URL not configured")
		}
		payload := map[string]string{
			"tg_id":       strconv.FormatInt(c.Sender().ID, 10),
			"tg_username": "@" + c.Sender().Username,
		}
		data, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", cfg.Backend.URL+"/api/auth/save_tg_user", bytes.NewReader(data))
		if err != nil {
			return c.Send("request error")
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(cfg.Backend.ClientID, cfg.Backend.AccessToken)
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return c.Send("backend unavailable")
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			return c.Send("registration failed")
		}
		return c.Send("Ты зарегистрирован. Для повышения роли, обратись к куратору")
	})
}
