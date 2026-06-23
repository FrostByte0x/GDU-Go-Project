package controllers_test

import (
	"testing"
	"wacdo-backend/controllers"
	"wacdo-backend/models"

	"github.com/maxatome/go-testdeep/td"
)

// createTestOrder creates an order containing one product.
func createTestOrder(t *testing.T) *models.Order {
	t.Helper()
	product := createTestProduct(t)
	input := models.OrderInput{
		Products: []models.OrderProduct{
			{ProductID: uint(product.ID), Quantity: 2},
		},
		Menus: []models.OrderMenu{},
	}
	order, err := controllers.CreateOrder(testDB, &input)
	td.CmpNoError(t, err)
	return order
}

// createTestOrderWithMenu creates an order containing one menu.
func createTestOrderWithMenu(t *testing.T) *models.Order {
	t.Helper()
	menu := createTestMenu(t)
	input := models.OrderInput{
		Products: []models.OrderProduct{},
		Menus: []models.OrderMenu{
			{MenuID: menu.ID, Quantity: 1},
		},
	}
	order, err := controllers.CreateOrder(testDB, &input)
	td.CmpNoError(t, err)
	return order
}

// createTestUnavailableMenu creates a menu marked as unavailable.
func createTestUnavailableMenu(t *testing.T) *models.Menu {
	t.Helper()
	product := createTestProduct(t)
	menuProduct := models.MenuProduct{
		Name:      product.Name,
		Quantity:  1,
		ProductID: uint(product.ID),
	}
	testMenu := models.Menu{
		Name:      "Unavailable Menu",
		Price:     product.UnitPrice,
		Products:  []models.MenuProduct{menuProduct},
		Available: false,
	}
	td.CmpNoError(t, controllers.CreateMenu(testDB, &testMenu))
	td.Cmp(t, testMenu.ID, td.NotZero())
	return &testMenu
}

func TestCreateOrder(t *testing.T) {
	t.Run("with products only", func(t *testing.T) {
		product := createTestProduct(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: uint(product.ID), Quantity: 3},
			},
			Menus: []models.OrderMenu{},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNoError(t, err)
		td.Cmp(t, order.ID, td.NotZero())
		td.Cmp(t, order.Price, td.Gt(0.0))
		td.Cmp(t, order.Price, product.UnitPrice*3)
		// State defaults to Created
		td.Cmp(t, order.State, models.Created)
		// Snapshot: name and price are filled in
		td.Cmp(t, order.Products[0].Name, product.Name)
		td.Cmp(t, order.Products[0].UnitPrice, product.UnitPrice)
	})

	t.Run("with menus only", func(t *testing.T) {
		menu := createTestMenu(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{},
			Menus: []models.OrderMenu{
				{MenuID: menu.ID, Quantity: 2},
			},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNoError(t, err)
		td.Cmp(t, order.ID, td.NotZero())
		td.Cmp(t, order.Price, menu.Price*2)
		td.Cmp(t, order.Menus[0].Name, menu.Name)
		td.Cmp(t, order.Menus[0].UnitPrice, menu.Price)
	})

	t.Run("with products and menus", func(t *testing.T) {
		product := createTestProduct(t)
		menu := createTestMenu(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: uint(product.ID), Quantity: 1},
			},
			Menus: []models.OrderMenu{
				{MenuID: menu.ID, Quantity: 1},
			},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNoError(t, err)
		td.Cmp(t, order.Price, product.UnitPrice+menu.Price)
	})

	t.Run("with unavailable product", func(t *testing.T) {
		unavailableProduct := createTestUnavailableProduct(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: uint(unavailableProduct.ID), Quantity: 1},
			},
			Menus: []models.OrderMenu{},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNil(t, order)
		td.CmpError(t, err)
	})

	t.Run("with unavailable menu returns error", func(t *testing.T) {
		menu := createTestUnavailableMenu(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{},
			Menus: []models.OrderMenu{
				{MenuID: menu.ID, Quantity: 1},
			},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNil(t, order)
		td.CmpError(t, err)
	})

	t.Run("price sums multiple products correctly", func(t *testing.T) {
		p1 := createTestProduct(t)
		p2 := createTestProduct(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: uint(p1.ID), Quantity: 2},
				{ProductID: uint(p2.ID), Quantity: 3},
			},
			Menus: []models.OrderMenu{},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNoError(t, err)
		td.Cmp(t, order.Price, p1.UnitPrice*2+p2.UnitPrice*3)
	})

	t.Run("price snapshot is not affected by later product price changes", func(t *testing.T) {
		product := createTestProduct(t)
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: uint(product.ID), Quantity: 2},
			},
			Menus: []models.OrderMenu{},
		}
		order, err := controllers.CreateOrder(testDB, &input)
		td.CmpNoError(t, err)
		originalPrice := order.Price
		// Change the product price after order creation
		_, err = controllers.UpdateProduct(testDB, product.ID, map[string]any{"unit_price": 999.99})
		td.CmpNoError(t, err)
		// Re-fetch the order: its price and snapshot unit price must be unchanged
		refetched, err := controllers.GetOrder(testDB, int(order.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, refetched.Price, originalPrice)
		td.Cmp(t, refetched.Products[0].UnitPrice, product.UnitPrice)
	})

	t.Run("empty order is rejected", func(t *testing.T) {
		input := models.OrderInput{
			Products: []models.OrderProduct{},
			Menus:    []models.OrderMenu{},
		}
		_, err := controllers.CreateOrder(testDB, &input)
		td.CmpError(t, err)
	})

	t.Run("non-existent product ID returns error", func(t *testing.T) {
		input := models.OrderInput{
			Products: []models.OrderProduct{
				{ProductID: 99999, Quantity: 1},
			},
			Menus: []models.OrderMenu{},
		}
		_, err := controllers.CreateOrder(testDB, &input)
		td.CmpError(t, err)
	})

	t.Run("non-existent menu ID returns error", func(t *testing.T) {
		input := models.OrderInput{
			Products: []models.OrderProduct{},
			Menus: []models.OrderMenu{
				{MenuID: 99999, Quantity: 1},
			},
		}
		_, err := controllers.CreateOrder(testDB, &input)
		td.CmpError(t, err)
	})
}

func TestDeleteOrder(t *testing.T) {
	t.Run("deletes order in Created state", func(t *testing.T) {
		order := createTestOrder(t)
		td.Cmp(t, order.State, models.Created)
		err := controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpNoError(t, err)
		// Confirm the order is gone
		_, err = controllers.GetOrder(testDB, int(order.ID))
		td.CmpError(t, err)
	})

	t.Run("cannot delete Validated order", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		err = controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpError(t, err)
	})

	t.Run("cannot delete Ready order", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Ready)
		td.CmpNoError(t, err)
		err = controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpError(t, err)
	})

	t.Run("cannot delete Delivered order", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Delivered)
		td.CmpNoError(t, err)
		err = controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpError(t, err)
	})

	t.Run("non-existent order returns error", func(t *testing.T) {
		err := controllers.DeleteOrder(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("deleted order is not returned by GetOrder", func(t *testing.T) {
		order := createTestOrder(t)
		err := controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpNoError(t, err)
		_, err = controllers.GetOrder(testDB, int(order.ID))
		td.CmpError(t, err)
	})

	t.Run("deleted order is not returned by GetOrders", func(t *testing.T) {
		order := createTestOrder(t)
		err := controllers.DeleteOrder(testDB, int(order.ID))
		td.CmpNoError(t, err)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{})
		td.CmpNoError(t, err)
		for _, o := range orders {
			td.Cmp(t, o.ID, td.Not(order.ID))
		}
	})
}

func TestGetOrders(t *testing.T) {
	t.Run("returns all orders when no filter", func(t *testing.T) {
		createTestOrder(t)
		createTestOrder(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(2))
	})

	t.Run("filters by Created state", func(t *testing.T) {
		createTestOrder(t)
		state := models.Created
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{State: &state})
		td.CmpNoError(t, err)
		for _, o := range orders {
			td.Cmp(t, o.State, models.Created)
		}
	})

	t.Run("filters by Validated state", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		state := models.Validated
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{State: &state})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		for _, o := range orders {
			td.Cmp(t, o.State, models.Validated)
		}
	})

	t.Run("filters by Ready state", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Ready)
		td.CmpNoError(t, err)
		state := models.Ready
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{State: &state})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		for _, o := range orders {
			td.Cmp(t, o.State, models.Ready)
		}
	})

	t.Run("filters by Delivered state", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Delivered)
		td.CmpNoError(t, err)
		state := models.Delivered
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{State: &state})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		for _, o := range orders {
			td.Cmp(t, o.State, models.Delivered)
		}
	})

	t.Run("sort=asc succeeds and returns results", func(t *testing.T) {
		createTestOrder(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{Sort: "asc"})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
	})

	t.Run("sort=desc succeeds and returns results", func(t *testing.T) {
		createTestOrder(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{Sort: "desc"})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
	})

	t.Run("sort=asc and sort=desc return same count", func(t *testing.T) {
		asc, err := controllers.GetOrders(testDB, models.OrderFilter{Sort: "asc"})
		td.CmpNoError(t, err)
		desc, err := controllers.GetOrders(testDB, models.OrderFilter{Sort: "desc"})
		td.CmpNoError(t, err)
		td.Cmp(t, len(asc), len(desc))
	})

	t.Run("preloads products and menus", func(t *testing.T) {
		createTestOrder(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		for _, o := range orders {
			if len(o.Products) > 0 {
				td.Cmp(t, o.Products[0].Name, td.NotZero())
				return
			}
		}
	})

	t.Run("preloads menus on menu-only orders", func(t *testing.T) {
		createTestOrderWithMenu(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{})
		td.CmpNoError(t, err)
		for _, o := range orders {
			if len(o.Menus) > 0 {
				td.Cmp(t, o.Menus[0].Name, td.NotZero())
				return
			}
		}
	})

	t.Run("state filter combined with sort=desc", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		state := models.Validated
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{State: &state, Sort: "desc"})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		for _, o := range orders {
			td.Cmp(t, o.State, models.Validated)
		}
	})
}

func TestGetOrder(t *testing.T) {
	t.Run("returns existing order with items", func(t *testing.T) {
		created := createTestOrder(t)
		order, err := controllers.GetOrder(testDB, int(created.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, order.ID, created.ID)
		td.Cmp(t, order.Price, created.Price)
		// Products are preloaded
		td.Cmp(t, len(order.Products), td.Gte(1))
		td.Cmp(t, order.Products[0].Name, td.NotZero())
	})

	t.Run("menus are preloaded on menu-only order", func(t *testing.T) {
		created := createTestOrderWithMenu(t)
		order, err := controllers.GetOrder(testDB, int(created.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, len(order.Menus), td.Gte(1))
		td.Cmp(t, order.Menus[0].Name, td.NotZero())
		td.Cmp(t, order.Menus[0].UnitPrice, td.Gt(0.0))
	})

	t.Run("returned order state matches", func(t *testing.T) {
		created := createTestOrder(t)
		order, err := controllers.GetOrder(testDB, int(created.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, order.State, models.Created)
	})

	t.Run("non-existent ID returns error", func(t *testing.T) {
		_, err := controllers.GetOrder(testDB, 99999)
		td.CmpError(t, err)
	})

	t.Run("negative ID returns error", func(t *testing.T) {
		_, err := controllers.GetOrder(testDB, -1)
		td.CmpError(t, err)
	})
}

func TestUpdateOrderState(t *testing.T) {
	t.Run("Created to Validated", func(t *testing.T) {
		order := createTestOrder(t)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Validated)
	})

	t.Run("Validated to Ready", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Ready)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Ready)
	})

	t.Run("Ready to Delivered", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Ready)
		td.CmpNoError(t, err)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Delivered)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Delivered)
	})

	t.Run("full lifecycle Created -> Validated -> Ready -> Delivered", func(t *testing.T) {
		order := createTestOrder(t)
		td.Cmp(t, order.State, models.Created)

		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Validated)

		updated, err = controllers.UpdateOrderState(testDB, int(order.ID), models.Ready)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Ready)

		updated, err = controllers.UpdateOrderState(testDB, int(order.ID), models.Delivered)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Delivered)
	})

	t.Run("state update persists when re-fetched", func(t *testing.T) {
		order := createTestOrder(t)
		_, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		refetched, err := controllers.GetOrder(testDB, int(order.ID))
		td.CmpNoError(t, err)
		td.Cmp(t, refetched.State, models.Validated)
	})

	t.Run("updating to same state succeeds", func(t *testing.T) {
		order := createTestOrder(t)
		td.Cmp(t, order.State, models.Created)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Created)
		td.CmpNoError(t, err)
		td.Cmp(t, updated.State, models.Created)
	})

	t.Run("non-existent order returns error", func(t *testing.T) {
		_, err := controllers.UpdateOrderState(testDB, 99999, models.Validated)
		td.CmpError(t, err)
	})

	t.Run("preloads products on return", func(t *testing.T) {
		order := createTestOrder(t)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		td.Cmp(t, len(updated.Products), td.Gte(1))
	})

	t.Run("preloads menus on return", func(t *testing.T) {
		order := createTestOrderWithMenu(t)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		td.Cmp(t, len(updated.Menus), td.Gte(1))
	})
}
