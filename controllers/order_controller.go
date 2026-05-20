package controllers

import (
	"wacdo-backend/models"

	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB, order *models.Order) error {
	return db.Create(order).Error
}
