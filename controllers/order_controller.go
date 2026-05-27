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

// Creating an order is a multi step transactions:
// 1. For each ordered Products, load its snapshot in the database.
// 2. For each ordered Menu, load its snapshot in the database.
// 3. Calculate the price of all of these.
// 4. Insert the whole thing in the DB as a single transaction.
// 5. Return the corresponding value to the client.

// CreateOrder inserts an order in the database and updates the pointer Order.
// It will also automatically calculate the order price.
func CreateOrder(db *gorm.DB, inputOrder *models.OrderInput) (*models.Order, error) {
	// Ensure the order is not empty, we don't take empty orders
	if len(inputOrder.Menus) < 1 && len(inputOrder.Products) < 1 {
		return nil, fmt.Errorf("At least one Menu or Product must be present in the order.")
	}
	// Create the complete Order
	var order models.Order
	order.Products = inputOrder.Products
	order.Menus = inputOrder.Menus
	// Store the price of items in a float64 var that will be added to the Order
	var sumOfPrice float64
	// Snapshot products Name and Price
	for k, v := range order.Products {
		product, err := GetProductByID(db, int(v.ProductID))
		if err != nil {
			return nil, err
		}
		order.Products[k].Name = product.Name
		order.Products[k].UnitPrice = product.UnitPrice
		sumOfPrice += order.Products[k].UnitPrice * float64(order.Products[k].Quantity)
	}
	// Snapshot menus
	for k, v := range order.Menus {
		menu, err := GetMenu(db, int(v.MenuID))
		if err != nil {
			return nil, err
		}
		order.Menus[k].Name = menu.Name
		order.Menus[k].UnitPrice = menu.Price
		sumOfPrice += order.Menus[k].UnitPrice * float64(order.Menus[k].Quantity)
	}
	// Set the price of the order
	order.Price = sumOfPrice
	err := db.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func CreateOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inputOrder models.OrderInput
		if err := c.ShouldBindJSON(&inputOrder); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order payload."})
			return
		}
		order, err := CreateOrder(db, &inputOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusCreated, order)

	}
}

// DeleteOrder deletes an order from the database.
//
// Only an admin should be allowed to perform this operation.
//
// Orders that are not "Created" cannot be deleted. They must be switched back to Created by an admin.
func DeleteOrder(db *gorm.DB, id int) error {
	var order models.Order
	if err := db.First(&order, id).Error; err != nil {
		return err
	}
	// comment this to allow deletion of orders that are in a higher state than created
	if order.State != models.Created {
		return fmt.Errorf("Cannot delete order with ID %d because its current state is %s", order.ID, order.State)
	}
	return db.Delete(&order).Error
}

// DeleteOrderHandler will handle requests
func DeleteOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure the user has the administrator role
		// if c.GetString("role") != "administrator" {
		// 	c.JSON(http.StatusForbidden, gin.H{"error": "Operation is not allowed"})
		// 	return
		// }
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
			return
		}
		if err := DeleteOrder(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// GetOrders will load orders from the database
func GetOrders(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	if err := db.Preload("Products").Preload("Menus").Find(&orders).Error; err != nil {
		return orders, err
	}
	return orders, nil
}

// GetOrdersHandler handles http requests to get orders
func GetOrdersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := GetOrders(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

// GetOrder will return a single order by its ID
func GetOrder(db *gorm.DB, id int) (*models.Order, error) {
	var order models.Order
	err := db.First(&order, id).Error
	return &order, err
}

// GetOrderHandler is the http handler that receives request to get a single order by its ID
func GetOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idPram := c.Param("id")
		id, err := strconv.Atoi(idPram)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		order, err := GetOrder(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with ID %d not found", id)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, &order)
	}
}

func UpdateOrderState(db *gorm.DB, id int, state models.OrderState) error {
	var order models.Order
	update := make(map[string]models.OrderState)
	update["state"] = state
	if err := db.First(&order, id).Error; err != nil {
		return err
	}
	return db.Model(&order).Updates(update).Error

}

func UpdateOrderHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		var StateUpdate models.StateOrderUpdate
		if err := c.ShouldBindJSON(&StateUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		if err := UpdateOrderState(db, id, StateUpdate.State); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Order with ID %d not found", id)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
