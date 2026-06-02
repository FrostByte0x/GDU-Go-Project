package routes

import (
	"wacdo-backend/controllers"
	"wacdo-backend/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(db *gorm.DB, router *gin.Engine) {
	routes := router.Group("/orders")
	routes.Use(middlewares.Authenticate())
	{
		routes.POST("", controllers.CreateOrderHandler(db))
		routes.DELETE("/:id", controllers.DeleteOrderHandler(db))
		// routes.Use(middlewares.Authorize([]models.Role{models.Preparator})).GET("", controllers.GetOrdersHandler(db))
		routes.GET("", controllers.GetOrdersHandler(db))
		routes.GET("/:id", controllers.GetOrderHandler(db))
		routes.PUT("/:id", controllers.UpdateOrderStateHandler(db))
	}
}
