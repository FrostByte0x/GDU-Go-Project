package controllers

import (
	"net/http"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Database operation to create a product
func CreateProduct(db *gorm.DB, product *models.Product) error {
	result := db.Create(product)
	return result.Error
}

// http handler to receive create product
func CreateProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product"})
			return
		}
		if err := CreateProduct(db, &product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create product"})
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

// Get products from the backend
func Getproducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	result := db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil

}

// return the products to the caller
func GetProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := Getproducts(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving products"})
			return
		}
		if len(products) == 0 {
			c.JSON(http.StatusOK, gin.H{"Success": "No products found"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}
