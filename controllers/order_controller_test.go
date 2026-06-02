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
		// Price must equal unitPrice * quantity
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

	t.Run("preloads products and menus", func(t *testing.T) {
		createTestOrder(t)
		orders, err := controllers.GetOrders(testDB, models.OrderFilter{})
		td.CmpNoError(t, err)
		td.Cmp(t, len(orders), td.Gte(1))
		// The first order with products should have them preloaded
		for _, o := range orders {
			if len(o.Products) > 0 {
				td.Cmp(t, o.Products[0].Name, td.NotZero())
				return
			}
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

	t.Run("non-existent order returns error", func(t *testing.T) {
		_, err := controllers.UpdateOrderState(testDB, 99999, models.Validated)
		td.CmpError(t, err)
	})

	t.Run("preloads products and menus on return", func(t *testing.T) {
		order := createTestOrder(t)
		updated, err := controllers.UpdateOrderState(testDB, int(order.ID), models.Validated)
		td.CmpNoError(t, err)
		td.Cmp(t, len(updated.Products), td.Gte(1))
	})
}
