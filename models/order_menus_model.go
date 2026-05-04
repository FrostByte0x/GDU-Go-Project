package models

// OrderMenu is the menus in a single Order
type OrderMenu struct {
	OrderID   uint   `gorm:"primaryKey;index"`
	MenuID    uint   `gorm:"primaryKey"`
	Quantity  uint   `gorm:"not null"`
	Name      string `gorm:"size:32"`  // Must match Menu name constraints!
	UnitPrice uint   `gorm:"not null"` // Price of the menu in the order
}
