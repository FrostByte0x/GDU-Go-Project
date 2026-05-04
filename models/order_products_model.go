package models

// OrderProduct is the products in a given Order that are not a menu
type OrderProduct struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"index"`
	ProductID uint    `gorm:"index"`
	Quantity  int     `gorm:"not null"`
	Name      string  `gorm:"size:32"`
	UnitPrice float64 `gorm:"type:decimal(10,2);not null"`
}
