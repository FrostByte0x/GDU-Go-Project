// package routes handles the routing of requests to resources
package routes

import (
	"wacdo-backend/controllers"
	"wacdo-backend/middlewares"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMenuRoutes(db *gorm.DB, router *gin.Engine) {
	routes := router.Group("/menus")
	{
		// Anyone can see the menus and the products -> usage métier : borne de commandes
		routes.GET("/:id", controllers.GetMenuHandler(db))
		routes.GET("", controllers.GetMenusHandler(db))
		{
			administrator := routes.Group("")
			administrator.Use(
				middlewares.Authenticate(),
				middlewares.Authorize([]models.Role{models.Administrator}),
			)
			// Expose CRUD operations for admin
			administrator.POST("", controllers.CreateMenuHandler(db))
			administrator.PUT("/:id", controllers.UpdateMenuHandler(db))
			administrator.DELETE("/:id", controllers.DeleteMenuHandler(db))
		}

	}
}
