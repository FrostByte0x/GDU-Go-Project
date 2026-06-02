package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LocalHostOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Request received from", "IP", c.ClientIP())
		if c.ClientIP() != "127.0.0.1" && c.ClientIP() != "::1" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.Next()
	}
}
