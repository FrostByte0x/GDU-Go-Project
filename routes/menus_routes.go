// package routes handles the routing of requests to resources
package routes

import (
	"wacdo-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMenuRoutes(db *gorm.DB, router *gin.Engine) {
	routes := router.Group("/menus")
	// Register the routes
	{
		routes.GET("", controllers.GetMenusHandler(db))
		routes.POST("", controllers.CreateMenuHandler(db))
		routes.GET("/:id", controllers.GetMenuHandler(db))
		routes.PUT("/:id", controllers.UpdateMenuHandler(db))
		routes.DELETE("/:id", controllers.DeleteMenuHandler(db))
	}
}
