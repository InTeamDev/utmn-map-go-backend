package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/tgbot"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
	telegramBot *tgbot.TelegramBot
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *service.AuthService, telegramBot *tgbot.TelegramBot) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		telegramBot: telegramBot,
	}
}

// RegisterHandlers registers all authentication endpoints
func (h *AuthHandler) RegisterHandlers(router *gin.Engine) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/register/generate-hash", h.GenerateRegistrationHash)
		auth.POST("/register", h.Register)
		auth.POST("/login/generate-code", h.GenerateAuthCode)
		auth.POST("/login", h.Login)
	}
}

// GenerateRegistrationHash godoc
// @Summary Generate a registration hash for Telegram deep link
// @Description Generate a hash that can be used in a Telegram deep link to register a user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} RegistrationHashResponse
// @Router /api/auth/register/generate-hash [post]
func (h *AuthHandler) GenerateRegistrationHash(c *gin.Context) {
	// Generate registration hash
	regHash := h.authService.GenerateRegistrationHash()

	// Create deep link
	deepLink := fmt.Sprintf("https://t.me/%s?start=%s", h.telegramBot.GetBotUsername(), regHash.Hash)

	// Return response
	c.JSON(http.StatusOK, RegistrationHashResponse{
		Hash:      regHash.Hash,
		DeepLink:  deepLink,
		ExpiresAt: regHash.ExpiresAt,
	})
}

// Register godoc
// @Summary Register a user from Telegram
// @Description Register a user after validating a registration hash from Telegram
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration info"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Validate hash
	if err := h.authService.ValidateRegistrationHash(req.Hash); err != nil {
		var status int
		var message string

		switch err {
		case service.ErrHashNotFound:
			status = http.StatusNotFound
			message = "Invalid registration hash"
		case service.ErrHashExpired:
			status = http.StatusUnauthorized
			message = "Registration hash has expired"
		default:
			status = http.StatusInternalServerError
			message = "Failed to validate registration hash"
		}

		c.JSON(status, ErrorResponse{
			Error: message,
		})
		return
	}

	// At this point, we would normally create the user
	// But in our design, user creation is done by the Telegram bot

	// Consume hash
	if err := h.authService.ConsumeRegistrationHash(req.Hash); err != nil {
		// Just log error, don't fail the request
		fmt.Printf("Failed to consume registration hash: %v", err)
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Success: true,
		Message: "Registration successful",
	})
}

// GenerateAuthCode godoc
// @Summary Generate an authentication code
// @Description Generate a one-time code for authentication and send it via Telegram
// @Tags auth
// @Accept json
// @Produce json
// @Param request body GenerateAuthCodeRequest true "Auth code request"
// @Success 200 {object} GenerateAuthCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/auth/login/generate-code [post]
func (h *AuthHandler) GenerateAuthCode(c *gin.Context) {
	var req GenerateAuthCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Add @ prefix if missing
	username := req.Username
	if !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	// Generate auth code
	authCode, err := h.authService.GenerateAuthCode(c.Request.Context(), username)
	if err != nil {
		var status int
		var message string

		if errors.Is(err, service.ErrUserNotFound) {
			status = http.StatusNotFound
			message = "User not found"
		} else {
			status = http.StatusInternalServerError
			message = "Failed to generate authentication code"
		}

		c.JSON(status, ErrorResponse{
			Error: message,
		})
		return
	}

	// Send auth code via Telegram
	err = h.telegramBot.SendAuthCode(c.Request.Context(), username, authCode.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to send authentication code",
		})
		return
	}

	c.JSON(http.StatusOK, GenerateAuthCodeResponse{
		Success:   true,
		Message:   "Authentication code sent",
		ExpiresAt: authCode.ExpiresAt,
	})
}

// Login godoc
// @Summary Authenticate a user
// @Description Authenticate a user using a one-time code and generate a JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Add @ prefix if missing
	username := req.Username
	if !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	// Validate auth code
	err := h.authService.ConsumeAuthCode(context.Background(), username, req.Code)
	if err != nil {
		var status int
		var message string

		switch err {
		case service.ErrCodeNotFound:
			status = http.StatusNotFound
			message = "Authentication code not found"
		case service.ErrCodeExpired:
			status = http.StatusUnauthorized
			message = "Authentication code has expired"
		case service.ErrInvalidCode:
			status = http.StatusUnauthorized
			message = "Invalid authentication code"
		default:
			status = http.StatusInternalServerError
			message = "Failed to validate authentication code"
		}

		c.JSON(status, ErrorResponse{
			Error: message,
		})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateJWT(context.Background(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

// RegistrationHashResponse represents the response for GenerateRegistrationHash
type RegistrationHashResponse struct {
	Hash      string `json:"hash"`
	DeepLink  string `json:"deep_link"`
	ExpiresAt any    `json:"expires_at"`
}

// RegisterRequest represents the request for Register
type RegisterRequest struct {
	Hash string `json:"hash" binding:"required"`
}

// RegisterResponse represents the response for Register
type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GenerateAuthCodeRequest represents the request for GenerateAuthCode
type GenerateAuthCodeRequest struct {
	Username string `json:"username" binding:"required"`
}

// GenerateAuthCodeResponse represents the response for GenerateAuthCode
type GenerateAuthCodeResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ExpiresAt any    `json:"expires_at"`
}

// LoginRequest represents the request for Login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// LoginResponse represents the response for Login
type LoginResponse struct {
	Token string `json:"token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
