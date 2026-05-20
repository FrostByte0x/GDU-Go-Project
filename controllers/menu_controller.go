// package controllers handles the change of resource states.
package controllers

import (
	"net/http"
	"wacdo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMenu(db *gorm.DB, menu *models.Menu) error {
	return db.Create(menu).Error
}

func CreateMenuHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu received"})
		}

	}
}

func GetMenus(db *gorm.DB) ([]models.Menu, error) {
	var menus []models.Menu
	if err := db.Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

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
