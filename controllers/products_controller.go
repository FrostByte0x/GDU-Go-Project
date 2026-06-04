package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Database operation to create a product
func CreateProduct(db *gorm.DB, product *models.Product) error {
	return db.Create(product).Error
}

// http handler to receive create product
//
//	@tags			products
//	@summary		Create a product
//	@accept			json
//	@produce		json
//	@security		BearerAuth
//	@description	Create a product
//	@param			product	body		models.ReturnProduct	true	"Product payload"
//	@success		201		{object}	models.ReturnProduct
//	@failure		400		{object}	models.ErrorResponse	"Invalid request payload received"
//	@failure		500		{object}	models.ErrorResponse	"Internal error creating the product"
//	@router			/products [post]
func CreateProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.BindJSON(&product); err != nil {
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
func GetProducts(db *gorm.DB) ([]models.Product, error) {
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// http handler to return the products to the caller
//
//	@tags		products
//	@summary	Get all products
//	@produce	json
//	@success	200	{object}	[]models.ReturnProduct	"The list of products"
//	@router		/products [get]
//	@failure	500	{object}	models.ErrorResponse	"Error retrieving products"
func GetProductsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := GetProducts(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving products"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

// Get product by ID from the backend
func GetProductByID(db *gorm.DB, id int) (*models.Product, error) {
	var product models.Product
	// Find the product by ID
	result := db.First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

// http handler to get a product by its ID
//
//	@tags		products
//	@summary	Get a single product by its ID
//	@produce	json
//	@param		ID	path		int						true	"Product ID"
//	@success	200	{object}	models.ReturnProduct	"A product"
//	@router		/products/{ID} [get]
func GetProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		// Get the product using ID from the request param
		product, err := GetProductByID(db, id)
		if err != nil {
			// Ensure the error is record not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %d not found", id)})
				return
				// otherwise return generic error
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve requested product"})
				return
			}
		}
		// return the product to the caller
		c.JSON(http.StatusOK, product)
	}
}

// Delete Products
func DeleteProduct(db *gorm.DB, id int) error {
	var product *models.Product
	product, err := GetProductByID(db, id)
	if err != nil {
		return err
	}
	return db.Delete(&product).Error
}

// Delete product http handler
//
//	@router		/products/{ID} [delete]
//	@tags		products
//	@summary	Delete a product
//	@success	204
//	@param		ID	path	int	true	"Product ID"
//	@security	BearerAuth
//	@failure	400	{object}	models.ErrorResponse	"Bad request"
//	@failure	404	{object}	models.ErrorResponse	"Product does not exist"
//	@failure	500	{object}	models.ErrorResponse	"Internal server error"
func DeleteProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		err = DeleteProduct(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("product with ID %d not found.", id)})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
				return
			}
		}
		c.Status(http.StatusNoContent)
	}
}

// Update Products
// To update a product, we create an update struct for fields that are allowed to be updated by the client
// such as the price or description of the product
func UpdateProduct(db *gorm.DB, id int, update map[string]any) (*models.Product, error) {
	var product models.Product
	err := db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	if err := db.Model(&product).Updates(update).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProductHandler returns the http func that handles requests to update a product, using
// the UpdateProduct
//
//	@summary		Update a product
//	@description	Update a product properties
//	@param			ID		path	int						true	"Product ID"
//	@param			update	body	models.UpdateProducts	true	"fields to update"
//	@router			/products/{ID} [put]
//	@security		BearerAuth
//	@tags			products
//	@param			ID	path		int						true	"Product ID"
//	@success		200	{object}	models.ReturnProduct	"The updated product"
//	@produce		json
//	@failure		400	{object}	models.ErrorResponse	"Invalid request payload"
//	@failure		404	{object}	models.ErrorResponse	"Product not found"
//	@failure		500	{object}	models.ErrorResponse	"Error upading object"
//	@accept			json
func UpdateProductHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputUpdate models.UpdateProducts
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		if err := c.ShouldBindJSON(&inputUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}
		update := make(map[string]any)
		// Add relevant fields to the map before updating in database
		// The model struct uses pointers so we can check for their existance with nil
		if inputUpdate.Name != nil {
			update["name"] = *inputUpdate.Name
		}
		if inputUpdate.Type != nil {
			update["type"] = *inputUpdate.Type
		}
		if inputUpdate.UnitPrice != nil {
			update["unit_price"] = *inputUpdate.UnitPrice
		}
		if inputUpdate.Available != nil {
			update["available"] = *inputUpdate.Available
		}
		// ensure there are updates to perform
		if len(update) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No update received"})
			return
		}
		// Update the database information with the updated product
		product, err := UpdateProduct(db, id, update)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating product"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}
