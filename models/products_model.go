package models

import (
	"time"

	"gorm.io/gorm"
)

type Category string

const (
	Boisson        Category = "Boisson"
	Burger         Category = "Burger"
	Accompagnement Category = "Accompagnement"
)

type Product struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string   `gorm:"size:64;not null" json:"name" binding:"required"`
	UnitPrice float64  `gorm:"type:decimal(10,2);not null" json:"price" binding:"required"`
	Type      Category `gorm:"type:enum('Boisson','Burger','Accompagnement');not null" json:"type" binding:"required"`
	Available bool     `gorm:"not null;default:true" json:"available" binding:"required"`
}

// ReturnProduct is the product returned by the GET products API
type ReturnProduct struct {
	ID        int     `json:"ID"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"price"`
	Type      string  `json:"type"`
	Available bool    `json:"available"`
}

type UpdateProducts struct {
	Name      *string   `gorm:"size:64;not null" json:"name"`
	UnitPrice *float64  `gorm:"type:decimal(10,2);not null" json:"price"`
	Type      *Category `gorm:"type:enum('Boisson','Burger','Accompagnement');not null" json:"type"`
	Available *bool     `gorm:"not null;default:true" json:"available"`
}
