package main

import (
	"log"
	"log/slog"
	"os"
	"wacdo-backend/config"
	"wacdo-backend/models"
	"wacdo-backend/routes"

	_ "wacdo-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Wacdo Order System Backend
//	@version		1.0
//	@description	Le backend des commandes Wacdo gère pour chaque magasin: les commandes, les produits, les menus.
//	@contact.name	frostByte0x
//	@contact.email	frostByte0x@github.com
//	@contact.url	https://github.com/FrostByte0x/GDU-Go-Project

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.apiKey	BearerAuth
// @in							header
// @name						Authorization
// @description				Use a bearer token to authenticate. Get a token at /users/login
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
	// Add swagger and schema inspector. They do not receive the Cors / rate limit and security configuration.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.RegisterWatcher(db, router)

	// Add security configurations: rate limit, CORS policy, basic security
	router.Use(config.CORSMiddleware())
	router.Use(config.SecurityMiddleware())
	router.Use(config.RateLimit(100))
	// Register the routes
	routes.RegisterProductRoutes(db, router)
	routes.RegisterUserRoutes(db, router)
	routes.RegisterMenuRoutes(db, router)
	routes.RegisterOrderRoutes(db, router)
	// Start the web server
	slog.Info("Server started, listening on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
