// package middlewares implements authentication and authorization mecanisms.
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}
		bearer := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.ParseWithClaims(bearer, &JwtStruct{},
			func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("jwt: unexpected signing method")
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to authenticate"})
			return
		}
		claims, ok := token.Claims.(*JwtStruct)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to authenticate"})
			return
		}
		// Inject username and role in the context
		c.Set("role", claims.Role)
		c.Set("username", claims.Subject)
		c.Next()
	}
}

type JwtStruct struct {
	jwt.RegisteredClaims
	Role models.Role `json:"role"`
}
