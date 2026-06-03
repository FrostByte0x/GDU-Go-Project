package controllers_test

import (
	"testing"
	"wacdo-backend/controllers"
	"wacdo-backend/models"

	"github.com/maxatome/go-testdeep/td"
)

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
	t.Run("creates available product with correct fields", func(t *testing.T) {
		product := models.Product{
			Name:      "Burger",
			UnitPrice: 10,
			Type:      models.Burger,
			Available: true,
		}
		err := controllers.CreateProduct(testDB, &product)
		td.CmpNoError(t, err)
		td.Cmp(t, product.ID, td.NotZero())
		td.Cmp(t, product.Name, "Burger")
		td.Cmp(t, product.UnitPrice, 10.0)
		td.Cmp(t, product.Type, models.Burger)
		td.CmpTrue(t, product.Available)
	})

	t.Run("creates unavailable product", func(t *testing.T) {
		product := models.Product{
			Name:      "Fanta",
			UnitPrice: 3.5,
			Type:      models.Boisson,
			Available: false,
		}
		err := controllers.CreateProduct(testDB, &product)
		td.CmpNoError(t, err)
		td.Cmp(t, product.ID, td.NotZero())
		td.CmpFalse(t, product.Available)
	})

	t.Run("all three categories are accepted", func(t *testing.T) {
		for _, cat := range []models.Category{models.Burger, models.Boisson, models.Accompagnement} {
			p := models.Product{Name: "Item", UnitPrice: 5, Type: cat, Available: true}
			td.CmpNoError(t, controllers.CreateProduct(testDB, &p))
			td.Cmp(t, p.Type, cat)
			td.Cmp(t, p.ID, td.NotZero())
		}
	})

	t.Run("each created product gets a distinct ID", func(t *testing.T) {
		p1 := models.Product{Name: "P1", UnitPrice: 1, Type: models.Burger, Available: true}
		p2 := models.Product{Name: "P2", UnitPrice: 1, Type: models.Burger, Available: true}
		td.CmpNoError(t, controllers.CreateProduct(testDB, &p1))
		td.CmpNoError(t, controllers.CreateProduct(testDB, &p2))
		td.Cmp(t, p1.ID, td.Not(p2.ID))
	})

	t.Run("decimal unit price is stored accurately", func(t *testing.T) {
		product := models.Product{Name: "Fries", UnitPrice: 3.75, Type: models.Accompagnement, Available: true}
		td.CmpNoError(t, controllers.CreateProduct(testDB, &product))
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.Cmp(t, got.UnitPrice, td.Gt(3.74), td.Lt(3.76))
	})
}

func TestGetProduct(t *testing.T) {
	t.Run("returns correct product fields", func(t *testing.T) {
		product := createTestProduct(t)
		td.Cmp(t, product.ID, td.NotZero())
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.Cmp(t, got.ID, product.ID)
		td.Cmp(t, got.Name, "Burger")
		td.Cmp(t, got.UnitPrice, 10.0)
		td.Cmp(t, got.Type, models.Burger)
		td.CmpTrue(t, got.Available)
	})

	t.Run("returns unavailable product", func(t *testing.T) {
		product := createTestUnavailableProduct(t)
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.CmpFalse(t, got.Available)
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		_, err := controllers.GetProductByID(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("negative ID returns error", func(t *testing.T) {
		_, err := controllers.GetProductByID(testDB, -4)
		td.CmpError(t, err)
	})
}

func TestGetProducts(t *testing.T) {
	t.Run("returns at least all created products", func(t *testing.T) {
		createTestProduct(t)
		createTestProduct(t)
		products, err := controllers.GetProducts(testDB)
		td.CmpNoError(t, err)
		td.Cmp(t, len(products), td.Gte(2))
	})

	t.Run("includes unavailable products", func(t *testing.T) {
		unavailable := createTestUnavailableProduct(t)
		products, err := controllers.GetProducts(testDB)
		td.CmpNoError(t, err)
		var found bool
		for _, p := range products {
			if p.ID == unavailable.ID {
				td.CmpFalse(t, p.Available)
				found = true
				break
			}
		}
		td.CmpTrue(t, found)
	})

	t.Run("returned products have all fields populated", func(t *testing.T) {
		product := createTestProduct(t)
		products, err := controllers.GetProducts(testDB)
		td.CmpNoError(t, err)
		var found bool
		for _, p := range products {
			if p.ID == product.ID {
				td.Cmp(t, p.Name, product.Name)
				td.Cmp(t, p.UnitPrice, product.UnitPrice)
				td.Cmp(t, p.Type, product.Type)
				td.CmpTrue(t, p.Available)
				found = true
				break
			}
		}
		td.CmpTrue(t, found)
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("updates name", func(t *testing.T) {
		product := createTestProduct(t)
		updated, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{"name": "Cheese Burger"})
		td.CmpNoError(t, err)
		td.Cmp(t, updated.Name, td.Not("Burger"))
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.Cmp(t, got.Name, "Cheese Burger")
	})

	t.Run("updates price", func(t *testing.T) {
		product := createTestProduct(t)
		updated, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{"unit_price": 12.5})
		td.CmpNoError(t, err)
		td.Cmp(t, updated.UnitPrice, td.Gt(11.0), td.Lt(13.0))
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.Cmp(t, got.UnitPrice, td.Gt(11.0), td.Lt(13.0))
	})

	t.Run("toggles availability to false", func(t *testing.T) {
		product := createTestProduct(t)
		updated, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{"available": false})
		td.CmpNoError(t, err)
		td.CmpFalse(t, updated.Available)
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.CmpFalse(t, got.Available)
	})

	t.Run("toggles availability back to true", func(t *testing.T) {
		product := createTestUnavailableProduct(t)
		updated, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{"available": true})
		td.CmpNoError(t, err)
		td.CmpTrue(t, updated.Available)
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.CmpTrue(t, got.Available)
	})

	t.Run("updates multiple fields simultaneously", func(t *testing.T) {
		product := createTestProduct(t)
		updated, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{
			"unit_price": 12.0,
			"available":  false,
			"name":       "Cheese Burger",
		})
		td.CmpNoError(t, err)
		td.CmpFalse(t, updated.Available)
		td.Cmp(t, updated.UnitPrice, td.Gt(11.0), td.Lt(13.0))
		td.Cmp(t, updated.Name, td.Not("Burger"))
	})

	t.Run("update persists when re-fetched", func(t *testing.T) {
		product := createTestProduct(t)
		_, err := controllers.UpdateProduct(testDB, product.ID, map[string]any{"name": "Persisted"})
		td.CmpNoError(t, err)
		got, err := controllers.GetProductByID(testDB, product.ID)
		td.CmpNoError(t, err)
		td.Cmp(t, got.Name, "Persisted")
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		_, err := controllers.UpdateProduct(testDB, 99999, map[string]any{"name": "Ghost"})
		td.CmpError(t, err)
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("soft deletes a product", func(t *testing.T) {
		product := createTestProduct(t)
		err := controllers.DeleteProduct(testDB, product.ID)
		td.CmpNoError(t, err)
		_, err = controllers.GetProductByID(testDB, product.ID)
		td.CmpError(t, err)
	})

	t.Run("deleted product is absent from GetProducts", func(t *testing.T) {
		product := createTestProduct(t)
		err := controllers.DeleteProduct(testDB, product.ID)
		td.CmpNoError(t, err)
		products, err := controllers.GetProducts(testDB)
		td.CmpNoError(t, err)
		for _, p := range products {
			td.Cmp(t, p.ID, td.Not(product.ID))
		}
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		err := controllers.DeleteProduct(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("negative ID returns error", func(t *testing.T) {
		err := controllers.DeleteProduct(testDB, -1)
		td.CmpError(t, err)
	})
}
