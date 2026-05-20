package controllers_tests

import (
	"testing"
	"wacdo-backend/controllers"
	"wacdo-backend/models"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Neither of these tests will work for now since we can't use enums with SQLite
func TestCreateProduct(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})
	// Define test product
	var product = models.Product{
		Name:      "Produit_name",
		Available: true,
		Type:      models.Boisson, // SQLite cannot use enums so we need to either spin up a mariaDB instance or move to using text
		UnitPrice: 5,
	}
	err := controllers.CreateProduct(db, &product)
	assert.NoError(t, err)
}

// Neither of these tests will work for now since we can't use enums with SQLite
func TestGetProducts(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Product{})
	_, err := controllers.GetProducts(db)
	assert.NoError(t, err)
}
