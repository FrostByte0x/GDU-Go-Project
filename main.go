package main

import (
	"log"
	"log/slog"
	"os"
	"wacdo-backend/config"
	"wacdo-backend/models"
	"wacdo-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load .env variables using go dot env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// InitDB handles errors itself, not in main.
	db := config.InitDB()
	if db == nil {
		slog.Error("Unable to instance DB object")
		os.Exit(1)
	}
	slog.Info("Successfully connected to", "database", os.Getenv("DB_NAME"))
	err = db.AutoMigrate(models.Menu{},
		models.MenuProduct{},
		models.Order{},
		models.OrderMenu{},
		models.User{},
		models.Product{},
		models.OrderProduct{},
	)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Database setup is complete.")
	// Register gin components
	router := gin.Default()
	// Trust no proxy
	router.SetTrustedProxies(nil)
	// Register the routes
	routes.RegisterProductRoutes(db, router)
	routes.RegisterMenuRoutes(db, router)
	routes.RegisterWatcher(db, router)
	routes.RegisterOrderRoutes(db, router)
	// Start the web server
	slog.Info("Server started, listening on port 8080")
	router.Run(":8080")
}
