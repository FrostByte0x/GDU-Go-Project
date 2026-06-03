// package models creates the object models for database and client-server opeations
package models

import (
	"time"

	"gorm.io/gorm"
)

// Menu is the model for menus in the database
// https://gorm.io/docs/models.html
type Menu struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      // Gorm automatically handles these fields at create and update time
	UpdatedAt time.Time      // Gorm automatically handles these fields at create and update time
	DeletedAt gorm.DeletedAt `gorm:"index"`                                          // Allow soft delete and requests such as WHERE deleted_at IS NULL
	Name      string         `json:"name" gorm:"size:32; not null"`                  // size:32 will use varchar(32) instead of LONGTEXT
	Price     float64        `gorm:"type:decimal(10,2); not null" json:"unit_price"` // maximum 99999999.99
	// Foreign Keys
	Products  []MenuProduct `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE"`
	Available bool          `json:"available" gorm:"not null"`
}

type UpdateMenu struct {
	Name      *string        `json:"name"`
	Price     *float64       `json:"unit_price"`
	Products  *[]MenuProduct `json:"products" gorm:"foreignKey:MenuID"`
	Available *bool          `json:"available"`
}
