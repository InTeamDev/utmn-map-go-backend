package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
	authservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/middleware"
)

type AuthAPI struct {
	svc        *authservice.Service
	repo       *repository.InMemory
	secret     []byte
	basicID    string
	basicToken string
}

func NewAuthAPI(
	svc *authservice.Service,
	repo *repository.InMemory,
	secret []byte,
	basicID, basicToken string,
) *AuthAPI {
	return &AuthAPI{svc: svc, repo: repo, secret: secret, basicID: basicID, basicToken: basicToken}
}

func (a *AuthAPI) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/auth")
	api.POST("/send_code", a.SendCode)
	api.POST("/verify", a.Verify)
	api.POST("/refresh", a.Refresh)
	api.POST("/save_tg_user", middleware.BasicAuth(a.basicID, a.basicToken), a.SaveUser)
	api.POST("/logout", middleware.JWTAuth(middleware.JWTAuthConfig{Secret: a.secret, Repo: a.repo}), a.Logout)
}

func (a *AuthAPI) SaveUser(c *gin.Context) {
	var req struct {
		TGID       int64  `json:"tg_id"`
		TGUsername string `json:"tg_username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.TGUsername == "" {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := a.svc.RegisterUser(req.TGID, req.TGUsername); err != nil {
		slog.Error("Failed to register user", "error", err, "tg_username", req.TGUsername)
		if errors.Is(err, authservice.ErrConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (a *AuthAPI) SendCode(c *gin.Context) {
	var req struct {
		TGUsername string `json:"tg_username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.TGUsername == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	exp, err := a.svc.SendCode(req.TGUsername)
	if err != nil {
		switch err {
		case authservice.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case authservice.ErrTooMany:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please try again later"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"expires_at": exp.Format(time.RFC3339)})
}

func (a *AuthAPI) Verify(c *gin.Context) {
	var req struct {
		TGUsername string `json:"tg_username"`
		Code       string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	tokens, err := a.svc.VerifyCode(req.TGUsername, req.Code)
	if err != nil {
		slog.Error("Failed to verify code", "error", err, "tg_username", req.TGUsername)
		switch err {
		case authservice.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case authservice.ErrInvalidCode, authservice.ErrExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		case authservice.ErrTooMany:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many attempts"})
		case authservice.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}

func (a *AuthAPI) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	tokens, err := a.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}

func (a *AuthAPI) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if err := a.svc.Logout(token, req.RefreshToken); err != nil {
		if errors.Is(err, authservice.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
