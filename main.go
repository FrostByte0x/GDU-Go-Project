package main

import (
	"log"
	"log/slog"
	"os"
	"wacdo-backend/config"
	"wacdo-backend/models"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// InitDB handles errors itself, not in main. Need to ensure it's the right design?
	db := config.InitDB()
	slog.Info("Successfully connected to", "database", os.Getenv("DB_NAME"))
	err = db.AutoMigrate(models.Menu{}, models.MenuProduct{}, models.OrderMenu{}, models.Order{}, models.User{}, models.Product{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Database setup is complete.")

}
