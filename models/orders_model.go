package models

import (
	"time"

	"gorm.io/gorm"
)

// Orders can be marked as created, validated, ready or delivered
type OrderState string

const (
	Created   OrderState = "Created"   // At creation, needs payment confirmation at Reception users
	Validated OrderState = "Validated" // After payment, will be prepared by Preparation users
	Ready     OrderState = "Ready"     // Order is ready to be delivered to the customer
	Delivered OrderState = "Delivered" // The customer has picked the order, this is the final lifecycle step.
)

// Order is an order from a customer. It contains an array of Products and Menus.
// The server returns the ID and the price.
// The state is by default to created.
type Order struct {
	ID        uint `gorm:"primaryKey" json:"ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Price     float64    `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	State     OrderState `gorm:"type:enum('Created','Validated','Ready', 'Delivered');not null;index; default:'Created'" json:"state"`
	// Foreign Keys
	Products []OrderProduct `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" binding:"required" json:"products"`
	Menus    []OrderMenu    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" binding:"required" json:"menus"`
}

// OrderInput is the client sent Order payload.
//
// This is done to ensure the client cannot set values such as Price or State, which are handled by the source of truth (the backend).
type OrderInput struct {
	Products []OrderProduct `json:"products" binding:"required"` // The array is required, even if empty
	Menus    []OrderMenu    `json:"menus" binding:"required"`    // The array is required, even if empty
}

// Struct to receive Order Updates. We only allow updating the status.
// An order can be deleted, or the customer can create a new one if needed.
type OrderUpdate struct {
	ID    uint       `json:"ID"`
	State OrderState `json:"state"`
}

// The StateOrderUpdate struct is to receive state update for an order. Other field updates are not permitted.
// ShouldBindJson will fail if a wrong value is provided.
type StateOrderUpdate struct {
	State OrderState `json:"state" binding:"required"`
}

func (state OrderState) IsValid() bool {
	switch state {
	case Created, Validated, Ready, Delivered:
		return true
	}
	return false
}
