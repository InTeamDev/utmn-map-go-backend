package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/repository"
	"github.com/gin-gonic/gin"
)

type JWTAuthConfig struct {
	Secret []byte
	Repo   *repository.InMemory
}

func JWTAuth(cfg JWTAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(401)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			c.AbortWithStatus(401)
			return
		}
		signingInput := parts[0] + "." + parts[1]
		mac := hmac.New(sha256.New, cfg.Secret)
		mac.Write([]byte(signingInput))
		sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(sig), []byte(parts[2])) {
			c.AbortWithStatus(401)
			return
		}
		payloadData, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(payloadData, &payload); err != nil {
			c.AbortWithStatus(401)
			return
		}
		if exp, ok := payload["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
			c.AbortWithStatus(401)
			return
		}
		if jti, ok := payload["jti"].(string); ok {
			if cfg.Repo.IsBlacklisted(jti) {
				c.AbortWithStatus(401)
				return
			}
		}
		c.Set("payload", payload)
		c.Next()
	}
}
