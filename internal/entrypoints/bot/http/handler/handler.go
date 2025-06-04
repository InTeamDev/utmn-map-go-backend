package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/telebot.v3"
)

type BotHandler struct {
	bot *telebot.Bot
}

func NewBotHandler(bot *telebot.Bot) *BotHandler {
	return &BotHandler{bot: bot}
}

type MessageRequest struct {
	TelegramUserID int64  `json:"telegram_user_id"`
	Message        string `json:"message"`
}

func (h *BotHandler) SendMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if _, err := h.bot.Send(&telebot.User{ID: req.TelegramUserID}, req.Message); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "telegram send failed"})
		return
	}
	c.Status(http.StatusOK)
	c.Status(http.StatusOK)
}
