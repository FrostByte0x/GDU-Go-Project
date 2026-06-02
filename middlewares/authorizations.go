package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
)

// Accept one or more role to access the APIs
func Authorize(roles []models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimedRoleString, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token"})
			return
		}
		// type assertion to ensure the interface is correct
		claimedRole, ok := claimedRoleString.(models.Role)
		if !slices.Contains(roles, claimedRole) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%s is not allowed to access this content", string(claimedRole))})
			return
		}
		c.Next()
	}
}
