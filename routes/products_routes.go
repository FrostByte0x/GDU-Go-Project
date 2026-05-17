package routes

import (
	"wacdo-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProductRoutes(db *gorm.DB, router *gin.Engine) {
	productRoutes := router.Group("/products")

	// productRoutes.Use middlewares for later

	{
		productRoutes.POST("", controllers.CreateProductHandler(db))
		productRoutes.GET("", controllers.GetProductsHandler(db))
	}
}
