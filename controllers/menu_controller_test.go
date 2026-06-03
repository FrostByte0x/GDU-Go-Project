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
	product := createTestProduct(t)
	menuProduct := models.MenuProduct{
		Name:      product.Name,
		Quantity:  2,
		ProductID: uint(product.ID),
	}
	testMenu := models.Menu{
		Name:      "Test Menu",
		Price:     product.UnitPrice * 2,
		Products:  []models.MenuProduct{menuProduct},
		Available: true,
	}
	err := controllers.CreateMenu(testDB, &testMenu)
	if err != nil {
		log.Fatal(err)
	}
	td.Cmp(t, testMenu.ID, td.NotZero())
	return &testMenu
}

func TestCreateMenu(t *testing.T) {
	t.Run("snapshots product name on creation", func(t *testing.T) {
		product := createTestProduct(t)
		menu := models.Menu{
			Name:  "Snap Menu",
			Price: product.UnitPrice,
			Products: []models.MenuProduct{
				{ProductID: uint(product.ID), Quantity: 1},
			},
			Available: true,
		}
		err := controllers.CreateMenu(testDB, &menu)
		td.CmpNoError(t, err)
		td.Cmp(t, menu.ID, td.NotZero())
		td.Cmp(t, menu.Products[0].Name, product.Name)
	})

	t.Run("creates menu with multiple products", func(t *testing.T) {
		p1 := createTestProduct(t)
		p2 := createTestProduct(t)
		menu := models.Menu{
			Name:  "Combo",
			Price: p1.UnitPrice + p2.UnitPrice,
			Products: []models.MenuProduct{
				{ProductID: uint(p1.ID), Quantity: 1},
				{ProductID: uint(p2.ID), Quantity: 2},
			},
			Available: true,
		}
		err := controllers.CreateMenu(testDB, &menu)
		td.CmpNoError(t, err)
		td.Cmp(t, len(menu.Products), 2)
		td.Cmp(t, menu.Products[0].Name, p1.Name)
		td.Cmp(t, menu.Products[1].Name, p2.Name)
	})

	t.Run("non-existent product ID returns error", func(t *testing.T) {
		menu := models.Menu{
			Name:  "Bad Menu",
			Price: 10,
			Products: []models.MenuProduct{
				{ProductID: 99999, Quantity: 1},
			},
		}
		err := controllers.CreateMenu(testDB, &menu)
		td.CmpError(t, err)
	})

	t.Run("price and name are stored correctly", func(t *testing.T) {
		product := createTestProduct(t)
		menu := models.Menu{
			Name:      "Stored Menu",
			Price:     42.50,
			Products:  []models.MenuProduct{{ProductID: uint(product.ID), Quantity: 1}},
			Available: true,
		}
		td.CmpNoError(t, controllers.CreateMenu(testDB, &menu))
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, got.Name, "Stored Menu")
		td.Cmp(t, got.Price, td.Gt(42.49), td.Lt(42.51))
	})

	t.Run("each created menu gets a distinct ID", func(t *testing.T) {
		m1 := createTestMenu(t)
		m2 := createTestMenu(t)
		td.Cmp(t, m1.ID, td.Not(m2.ID))
	})
}

func TestGetMenu(t *testing.T) {
	t.Run("returns existing menu", func(t *testing.T) {
		menu := createTestMenu(t)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, got, td.NotNil())
		td.Cmp(t, got.ID, menu.ID)
		td.Cmp(t, got.Price, menu.Price)
		td.Cmp(t, got.Name, menu.Name)
	})

	t.Run("preloads products", func(t *testing.T) {
		menu := createTestMenu(t)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, len(got.Products), td.Gte(1))
		td.Cmp(t, got.Products[0].Name, td.NotZero())
		td.Cmp(t, got.Products[0].Quantity, td.Gt(0))
	})

	t.Run("available flag is returned", func(t *testing.T) {
		menu := createTestMenu(t)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.CmpTrue(t, got.Available)
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		_, err := controllers.GetMenu(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("negative ID returns error", func(t *testing.T) {
		_, err := controllers.GetMenu(testDB, -4)
		td.CmpError(t, err)
	})
}

func TestGetMenus(t *testing.T) {
	t.Run("returns at least all created menus", func(t *testing.T) {
		createTestMenu(t)
		createTestMenu(t)
		menus, err := controllers.GetMenus(testDB)
		td.CmpNoError(t, err)
		td.Cmp(t, len(menus), td.Gte(2))
	})

	t.Run("preloads products for each menu", func(t *testing.T) {
		createTestMenu(t)
		menus, err := controllers.GetMenus(testDB)
		td.CmpNoError(t, err)
		for _, m := range menus {
			if len(m.Products) > 0 {
				td.Cmp(t, m.Products[0].Name, td.NotZero())
				return
			}
		}
	})

	t.Run("returned menus have all fields populated", func(t *testing.T) {
		menu := createTestMenu(t)
		menus, err := controllers.GetMenus(testDB)
		td.CmpNoError(t, err)
		var found bool
		for _, m := range menus {
			if m.ID == menu.ID {
				td.Cmp(t, m.Name, menu.Name)
				td.Cmp(t, m.Price, menu.Price)
				found = true
				break
			}
		}
		td.CmpTrue(t, found)
	})
}

func TestDeleteMenu(t *testing.T) {
	t.Run("soft deletes a menu", func(t *testing.T) {
		menu := createTestMenu(t)
		err := controllers.DeleteMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		_, err = controllers.GetMenu(testDB, int(menu.ID))
		td.CmpError(t, err)
	})

	t.Run("deleted menu is absent from GetMenus", func(t *testing.T) {
		menu := createTestMenu(t)
		err := controllers.DeleteMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		menus, err := controllers.GetMenus(testDB)
		td.CmpNoError(t, err)
		for _, m := range menus {
			td.Cmp(t, m.ID, td.Not(menu.ID))
		}
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		err := controllers.DeleteMenu(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("negative ID returns error", func(t *testing.T) {
		err := controllers.DeleteMenu(testDB, -1)
		td.CmpError(t, err)
	})
}

func TestUpdateMenu(t *testing.T) {
	t.Run("updates name and it persists", func(t *testing.T) {
		menu := createTestMenu(t)
		_, err := controllers.UpdateMenu(testDB, int(menu.ID), map[string]any{"name": "Renamed Menu"})
		td.CmpNoError(t, err)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, got.Name, "Renamed Menu")
	})

	t.Run("updates price and it persists", func(t *testing.T) {
		menu := createTestMenu(t)
		_, err := controllers.UpdateMenu(testDB, int(menu.ID), map[string]any{"price": 99.99})
		td.CmpNoError(t, err)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, got.Price, td.Gt(99.98), td.Lt(100.0))
	})

	t.Run("toggles availability to false", func(t *testing.T) {
		menu := createTestMenu(t)
		_, err := controllers.UpdateMenu(testDB, int(menu.ID), map[string]any{"available": false})
		td.CmpNoError(t, err)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.CmpFalse(t, got.Available)
	})

	t.Run("toggles availability back to true", func(t *testing.T) {
		// Start with available=false by creating an unavailable menu directly
		product := createTestProduct(t)
		unavailableMenu := models.Menu{
			Name:      "Temp Unavailable",
			Price:     5,
			Products:  []models.MenuProduct{{ProductID: uint(product.ID), Quantity: 1}},
			Available: false,
		}
		td.CmpNoError(t, controllers.CreateMenu(testDB, &unavailableMenu))
		_, err := controllers.UpdateMenu(testDB, int(unavailableMenu.ID), map[string]any{"available": true})
		td.CmpNoError(t, err)
		got, err := controllers.GetMenu(testDB, int(unavailableMenu.ID))
		td.CmpNoError(t, err)
		td.CmpTrue(t, got.Available)
	})

	t.Run("updates multiple fields simultaneously", func(t *testing.T) {
		menu := createTestMenu(t)
		_, err := controllers.UpdateMenu(testDB, int(menu.ID), map[string]any{
			"name":      "Multi Update",
			"price":     55.0,
			"available": false,
		})
		td.CmpNoError(t, err)
		got, err := controllers.GetMenu(testDB, int(menu.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, got.Name, "Multi Update")
		td.Cmp(t, got.Price, td.Gt(54.9), td.Lt(55.1))
		td.CmpFalse(t, got.Available)
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		_, err := controllers.UpdateMenu(testDB, 99999, map[string]any{"name": "Ghost"})
		td.CmpError(t, err)
	})
}
