package main

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"wacdo-backend/config"
	"wacdo-backend/controllers"
	"wacdo-backend/models"
	"wacdo-backend/routes"

	_ "wacdo-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		// We use .env in dev, but variables are injected through the run time in production.
		slog.Info("Error loading .env file, if you are running this as a container, this is expected")
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
	// Bootstrap admin user
	//
	const adminuser string = "admin2"
	if _, err := controllers.GetUserByUsername(db, adminuser); err != nil {
		slog.Info("Admin1 not found, processing bootstrap")
		// only if not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			password, err := bcrypt.GenerateFromPassword([]byte(adminuser), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("error generating bootstrap user password: %s", err.Error())
			}
			err = controllers.CreateUser(db, &models.User{
				Username: adminuser,
				Password: string(password),
			})
			if err != nil {
				log.Fatalf("error creating bootstrap user password: %s", err.Error())
			}
			_, err = controllers.UpdateUserRole(db, adminuser, models.Administrator)
			if err != nil {
				log.Fatalf("error updating bootstrap user role: %s", err.Error())
			}
		} else {
			log.Fatalf("unable to check existing admin user: %s", err.Error())
		}
	} else {
		slog.Info("Bootstrap skipped because admin user already exists")
	}
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
