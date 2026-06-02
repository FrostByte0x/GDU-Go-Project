package controllers_test

import (
	"testing"
	"wacdo-backend/controllers"
	"wacdo-backend/models"

	"github.com/maxatome/go-testdeep/td"
)

// Helper function to create a default, testing product
func createTestProduct(t *testing.T) models.Product {
	t.Helper()
	TestProduct := models.Product{
		Name:      "Burger",
		UnitPrice: 10,
		Type:      models.Burger,
		Available: true,
	}
	td.CmpNoError(t, controllers.CreateProduct(testDB, &TestProduct))
	return TestProduct
}
func createTestUnavailableProduct(t *testing.T) models.Product {
	t.Helper()
	TestProduct := models.Product{
		Name:      "Boisson non disponible",
		UnitPrice: 3.25,
		Type:      models.Boisson,
		Available: false,
	}
	td.CmpNoError(t, controllers.CreateProduct(testDB, &TestProduct))
	return TestProduct
}
func TestCreateProduct(t *testing.T) {
	TestProduct := models.Product{
		Name:      "Burger",
		UnitPrice: 10,
		Type:      models.Burger,
		Available: true,
	}
	err := controllers.CreateProduct(testDB, &TestProduct)
	// Ensure error is nil
	td.CmpNoError(t, err)
	// Ensure the values are correctly set
	td.Cmp(t, TestProduct.Name, "Burger")
	// Ensure Gorm correctly updated the system ID
	td.Cmp(t, TestProduct.ID, td.NotZero())
	// Ensure the bool worked
	td.CmpTrue(t, TestProduct.Available)
}

func TestUpdateProduct(t *testing.T) {
	product := createTestProduct(t)
	update := make(map[string]any)
	update["unit_price"] = 12
	update["available"] = false
	update["name"] = "Cheese Burger"
	updatedProduct, err := controllers.UpdateProduct(testDB, product.ID, update)
	td.CmpNoError(t, err)
	// Ensure the availability is correctly changed
	td.CmpFalse(t, updatedProduct.Available)
	// Ensure the price is updated to 12
	td.Cmp(t, updatedProduct.UnitPrice, td.Gt(11.0), td.Lt(13.0))
	// Ensure the name is changed
	td.Cmp(t, updatedProduct.Name, td.Not("Burger"))
}

// Includes test for fake IDs
func TestGetProduct(t *testing.T) {
	product := createTestProduct(t)
	// Ensure the ID is updated at creation
	td.Cmp(t, product.ID, td.NotZero())
	TestProduct, err := controllers.GetProductByID(testDB, product.ID)
	td.CmpNoError(t, err)
	// Ensure the values are correctly set
	td.Cmp(t, TestProduct.Name, "Burger")
	// Ensure Gorm correctly updated the system ID
	td.Cmp(t, TestProduct.ID, td.NotZero())
	// Ensure the bool worked
	td.CmpTrue(t, TestProduct.Available)
	// Test fake ID
	_, err = controllers.GetProductByID(testDB, 9999)
	td.CmpError(t, err)
	// Test invalid ID
	_, err = controllers.GetProductByID(testDB, -4)
	td.CmpError(t, err)
}
