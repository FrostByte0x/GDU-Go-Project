package models

// MenuProduct is the products in a given Menu
type MenuProduct struct {
	MenuID    uint   `gorm:"primaryKey;index"`
	ProductId uint   `gorm:"primaryKey"`
	Quantity  int    `gorm:"not null"`
	Name      string `gorm:"size:64;not null"`
}
