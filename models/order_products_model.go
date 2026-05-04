package models

// OrderProduct is the products in a given Order that are not a menu
type OrderProduct struct {
	OrderID   uint   `gorm:"primaryKey;index"`
	ProductID uint   `gorm:"primaryKey"`
	Quantity  int    `gorm:"not null"`
	Name      string `gorm:"size:32"`
	Price     int    `gorm:"not null"`
}
