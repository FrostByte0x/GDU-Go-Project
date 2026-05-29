package models

// OrderProduct is the products in a given Order that are not a menu
type OrderProduct struct {
	ID        uint    `gorm:"primaryKey" json:"ID"`
	OrderID   uint    `gorm:"index"`
	ProductID uint    `gorm:"index" json:"product_id"`
	Quantity  uint    `gorm:"not null" json:"quantity"`
	Name      string  `gorm:"size:32" json:"name"`
	UnitPrice float64 `gorm:"type:decimal(10,2);not null" json:"unit_price"`
}
