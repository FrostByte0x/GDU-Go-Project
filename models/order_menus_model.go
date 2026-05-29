package models

// OrderMenu is the menus in a single Order
type OrderMenu struct {
	ID        uint    `gorm:"primaryKey" json:"ID"`
	OrderID   uint    `gorm:"index"`
	MenuID    uint    `gorm:"index" json:"menu_id"`
	Quantity  uint    `gorm:"not null" json:"quantity"`
	Name      string  `gorm:"size:32" json:"name"`                      // Must match Menu name constraints!
	UnitPrice float64 `gorm:"not null;decimal(10,2)" json:"unit_price"` // Price of the menu in the order
}
