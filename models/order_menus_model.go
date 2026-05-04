package models

// OrderMenu is the menus in a single Order
type OrderMenu struct {
	ID        uint   `gorm:"primaryKey"`
	OrderID   uint   `gorm:"index"`
	MenuID    uint   `gorm:"index"`
	Menu      Menu   `gorm:"foreignKey:MenuID"`
	Quantity  uint   `gorm:"not null"`
	Name      string `gorm:"size:32"`  // Must match Menu name constraints!
	UnitPrice uint   `gorm:"not null"` // Price of the menu in the order
}
