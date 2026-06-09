package routes

import (
	"wacdo-backend/controllers"
	"wacdo-backend/middlewares"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(db *gorm.DB, router *gin.Engine) {
	UserRoutes := router.Group("/users")
	UserRoutes.Use(middlewares.LocalHostOnly())

	{
		UserRoutes.POST("/register", controllers.Register(db))
		// UserRoutes.POST("", controllers.CreateUserHandler(db))
		UserRoutes.POST("/login", controllers.Login(db))
		// Only admin can update users
		adminRoute := UserRoutes.Group("").Use(middlewares.Authenticate())
		adminRoute.Use(middlewares.Authorize([]models.Role{models.Administrator}))
		adminRoute.PUT("/:username/role", controllers.UpdateUserRoleHandler(db))
	}
}
