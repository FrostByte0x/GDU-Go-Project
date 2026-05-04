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
	Name      string   `gorm:"size:64;not null"`
	Price     int      `gorm:"type:decimal(10,2);not null"`
	Type      Category `gorm:"type:enum('Boisson','Burger','Accompagnement');not null"`
	Available bool     `gorm:"not null;default:true"`
}
