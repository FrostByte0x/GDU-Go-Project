package models

// MenuProduct is the products in a given Menu
type MenuProduct struct {
	ID        uint `gorm:"primaryKey"`
	MenuID    uint `gorm:"index"`
	ProductID uint
	Quantity  int    `gorm:"not null"`
	Name      string `gorm:"size:64;not null"`
}
