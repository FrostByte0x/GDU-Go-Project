package models

// MenuProduct is the products in a given Menu
type MenuProduct struct {
	ID        uint   `gorm:"primaryKey"`
	MenuID    uint   `gorm:"not null;index"`
	ProductID uint   `gorm:"not null" json:"product_id"`
	Quantity  int    `gorm:"not null"`
	Name      string `gorm:"size:64;not null"`
}

// MenuProductResponse is the list of products in a menu returned by the API to the client.
// This layer is only for returns
type MenuProductResponse struct {
	Name     string `json:"name"`
	Quantity uint   `json:"quantity"`
}

type MenuResponse struct {
	ID       uint                  `json:"id"`
	Name     string                `json:"name"`
	Price    uint                  `json:"price"`
	Products []MenuProductResponse `json:"products"`
}
