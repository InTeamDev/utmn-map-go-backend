package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
	authservice "github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/service"
	"github.com/InTeamDev/utmn-map-go-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	svc        *authservice.Service
	repo       *repository.InMemory
	secret     []byte
	basicID    string
	basicToken string
}

func NewAuthAPI(svc *authservice.Service, repo *repository.InMemory, secret []byte, basicID, basicToken string) *AuthAPI {
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
		c.Status(http.StatusBadRequest)
		return
	}
	if err := a.svc.RegisterUser(req.TGID, req.TGUsername); err != nil {
		if err == authservice.ErrConflict {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusInternalServerError)
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
			c.Status(http.StatusNotFound)
		case authservice.ErrTooMany:
			c.Status(http.StatusTooManyRequests)
		default:
			c.Status(http.StatusServiceUnavailable)
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
		c.Status(http.StatusBadRequest)
		return
	}
	access, refresh, err := a.svc.VerifyCode(req.TGUsername, req.Code)
	if err != nil {
		switch err {
		case authservice.ErrNotFound:
			c.Status(http.StatusNotFound)
		case authservice.ErrInvalidCode, authservice.ErrExpired:
			c.Status(http.StatusBadRequest)
		case authservice.ErrTooMany:
			c.Status(http.StatusTooManyRequests)
		case authservice.ErrUnauthorized:
			c.Status(http.StatusUnauthorized)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: refresh, HttpOnly: true, Path: "/"})
	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (a *AuthAPI) Refresh(c *gin.Context) {
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	access, refresh, err := a.svc.RefreshToken(cookie.Value)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: refresh, HttpOnly: true, Path: "/"})
	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (a *AuthAPI) Logout(c *gin.Context) {
	cookie, _ := c.Request.Cookie("refresh_token")
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	rt := ""
	if cookie != nil {
		rt = cookie.Value
		http.SetCookie(c.Writer, &http.Cookie{Name: "refresh_token", Value: "", Expires: time.Unix(0, 0), Path: "/"})
	}
	a.svc.Logout(token, rt)
	c.Status(http.StatusNoContent)
}
