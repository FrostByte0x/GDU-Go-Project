package routes

import (
	"wacdo-backend/controllers"
	"wacdo-backend/middlewares"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(db *gorm.DB, router *gin.Engine) {
	routes := router.Group("/orders")
	routes.Use(middlewares.Authenticate())
	{
		// Routes for preparator + reception
		preparator := routes.Group("")
		preparator.Use(
			middlewares.Authorize([]models.Role{models.Preparator, models.Reception}),
		)
		preparator.PUT("/:id", controllers.UpdateOrderStateHandler(db))
		preparator.GET("", controllers.GetOrdersHandler(db))
		preparator.GET("/:id", controllers.GetOrderHandler(db))
		preparator.POST("", controllers.CreateOrderHandler(db))
	}
	{
		administrator := routes.Group("").Use(middlewares.Authorize([]models.Role{models.Administrator}))
		administrator.DELETE("/:id", controllers.DeleteOrderHandler(db))

	}
}
