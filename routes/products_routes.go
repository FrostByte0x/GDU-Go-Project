package routes

import (
	"wacdo-backend/controllers"
	"wacdo-backend/middlewares"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProductRoutes(db *gorm.DB, router *gin.Engine) {
	productRoutes := router.Group("/products")
	{
		// Unauthenticated routes -> bornes de commandes
		productRoutes.GET("/:id", controllers.GetProductHandler(db))
		productRoutes.GET("", controllers.GetProductsHandler(db))
		{
			// Administrator only operations: manage products
			administrator := productRoutes.Group("").Use(
				middlewares.Authenticate(),
				middlewares.Authorize([]models.Role{models.Administrator}),
			)
			administrator.POST("", controllers.CreateProductHandler(db))
			administrator.PUT("/:id", controllers.UpdateProductHandler(db))
			administrator.DELETE("/:id", controllers.DeleteProductHandler(db))
		}
	}
}
