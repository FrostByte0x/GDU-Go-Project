package routes

import (
	"wacdo-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(db *gorm.DB, router *gin.Engine) {
	UserRoutes := router.Group("/users")
	// UserRoutes.Use(middlewares.LocalHostOnly()) .Use on User management routes

	{
		UserRoutes.POST("/register", controllers.Register(db))
		// UserRoutes.POST("", controllers.CreateUserHandler(db))
		UserRoutes.POST("/login", controllers.Login(db))
		UserRoutes.PUT("/:username/role", controllers.UpdateUserRoleHandler(db))
	}
}
