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
		// If the role is allowed, OR if the role is administrator, which has all permissions.
		if slices.Contains(roles, claimedRole) || claimedRole == models.Administrator {
			c.Next()
			return
		}
		// If not, we deny the request
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%s is not allowed to access this content", string(claimedRole))})
	}
}
