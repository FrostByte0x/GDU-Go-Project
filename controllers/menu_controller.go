// package controllers handles the change of resource states.
package controllers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateMenu creates the menu in the dabase. Returns an error or nil.
func CreateMenu(db *gorm.DB, menu *models.Menu) error {
	// For each product, snapshot the name
	for k, v := range menu.Products {
		product, err := GetProductByID(db, int(v.ProductID))
		if err != nil {
			return err
		}
		menu.Products[k].Name = product.Name
	}
	return db.Create(menu).Error
}

// CreateMenuHandler is the http handler that receives the Create Menu request
func CreateMenuHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu received"})
			return
		}
		if err := CreateMenu(db, &menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating menu"})
			return
		}
		c.JSON(http.StatusCreated, ToMenuResponse(menu))
	}
}

// GetMenus will return all menus
func GetMenus(db *gorm.DB) ([]models.Menu, error) {
	var menus []models.Menu
	if err := db.Preload("Products").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// GetMenusHandler is the http handler that will receive the request and return all menus
func GetMenusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := GetMenus(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error loading menus"})
			return
		}
		c.JSON(http.StatusOK, menus)
	}
}

// GetMenu will get a single menu by ID
func GetMenu(db *gorm.DB, id int) (*models.Menu, error) {
	var menu models.Menu
	if err := db.Preload("Products").First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetMenusHandler is the http handler that receives
func GetMenuHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		menu, err := GetMenu(db, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Menu with ID %d not found", id)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error retrieving menu"})
			return
		}
		slog.Info(fmt.Sprintf("There are %d products in the requested menu", len(menu.Products)))
		c.JSON(http.StatusOK, ToMenuResponse(*menu))
	}
}

// DeleteMenu will delete a menu from the database
func DeleteMenu(db *gorm.DB, id int) error {
	var menu models.Menu
	// Find the menu
	if err := db.First(&menu, id).Error; err != nil {
		return err
	}
	// Delete it
	return db.Delete(&menu).Error
}

// DeleteMenuHandler is the http handler that receives DELETE requests for a given menu
func DeleteMenuHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		if err := DeleteMenu(db, id); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Menu with ID %d not found.", id)})
				return
			}
		}
		c.Status(http.StatusNoContent)
	}
}

func UpdateMenu(db *gorm.DB, id int, update map[string]any) (*models.Menu, error) {
	var menu models.Menu
	if err := db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&menu).Updates(update).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func UpdateMenuHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateMenu models.UpdateMenu
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		if err := c.ShouldBindJSON(&updateMenu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		update := make(map[string]any)
		if updateMenu.Name != nil {
			update["name"] = *updateMenu.Name
		}
		if updateMenu.Price != nil {
			update["price"] = *updateMenu.Price
		}
		if updateMenu.Products != nil {
			update["products"] = *updateMenu.Products
		}
		if updateMenu.Available != nil {
			update["available"] = *updateMenu.Available
		}
		if len(update) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty update payload"})
			return
		}
		updatedMenu, err := UpdateMenu(db, id, update)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Menu with ID %d not found", id)})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating Menu"})
			return
		}
		c.JSON(http.StatusOK, ToMenuResponse(*updatedMenu))
	}
}

func ToMenuResponse(menu models.Menu) models.MenuResponse {
	products := make([]models.MenuProductResponse, len(menu.Products))
	for k, v := range menu.Products {
		products[k] = models.MenuProductResponse{
			Name:     v.Name,
			Quantity: uint(v.Quantity),
		}
	}
	return models.MenuResponse{
		ID:       menu.ID,
		Name:     menu.Name,
		Price:    uint(menu.Price),
		Products: products,
	}
}
