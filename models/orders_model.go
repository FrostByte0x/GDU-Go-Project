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

type Order struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Price     float64    `gorm:"type:decimal(10,2);not null"`
	State     OrderState `gorm:"type:enum('Created','Validated','Ready', 'Delivered');not null;index"`
	// Foreign Keys
	Products []OrderProduct `gorm:"foreignKey:OrderID"`
	Menus    []Menu         `gorm:"foreignKey:OrderID"`
}
