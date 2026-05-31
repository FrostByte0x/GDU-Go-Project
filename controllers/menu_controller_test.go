package controllers_test

import (
	"log"
	"testing"
	"wacdo-backend/controllers"
	"wacdo-backend/models"

	"github.com/maxatome/go-testdeep/td"
)

func createTestMenu(t *testing.T) *models.Menu {
	t.Helper()
	// Create a test product
	product := createTestProduct(t)
	menuProduct := models.MenuProduct{
		Name:      product.Name,
		Quantity:  2,
		ProductID: uint(product.ID),
	}
	testMenu := models.Menu{
		Name:     "Test Menu",
		Price:    product.UnitPrice * 2,
		Products: []models.MenuProduct{menuProduct},
	}
	err := controllers.CreateMenu(testDB, &testMenu)
	if err != nil {
		log.Fatal(err)
	}
	// Ensure we get a return ID in the Menu
	td.Cmp(t, testMenu.ID, td.NotZero())
	return &testMenu
}

func TestGetMenu(t *testing.T) {
	menu := createTestMenu(t)
	getMenu, err := controllers.GetMenu(testDB, int(menu.ID))
	td.CmpNoError(t, err)
	td.Cmp(t, getMenu, td.NotNil())
	// Test fake ID
	_, err = controllers.GetMenu(testDB, 9999)
	td.CmpError(t, err)
	_, err = controllers.GetMenu(testDB, -4)
	td.CmpError(t, err)
	td.Cmp(t, getMenu.Price, menu.Price)
}
