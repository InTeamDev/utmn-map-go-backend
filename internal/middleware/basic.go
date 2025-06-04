package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth(expectedID, expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Basic ") {
			c.AbortWithStatus(401)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 || parts[0] != expectedID || parts[1] != expectedToken {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
