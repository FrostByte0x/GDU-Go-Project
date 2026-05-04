package main

import (
	"log/slog"
	"os"
	"wacdo-backend/config"
	"wacdo-backend/models"
)

func main() {
	// InitDB handles errors itself, not in main. Need to ensure it's the right design?
	db := config.InitDB()
	slog.Info("Successfully connected to", "database", os.Getenv("DB_NAME"))
	db.AutoMigrate(models.Menu{}, models.MenuProduct{}, models.OrderMenu{}, models.Order{}, models.User{}, models.Product{})
}
