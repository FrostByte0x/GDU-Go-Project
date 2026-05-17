package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(db *gorm.DB, router *gin.Engine) {
	UserRoutes := router.Group("/users")

	{
		UserRoutes.POST("")
	}
}
