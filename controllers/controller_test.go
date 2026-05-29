package controllers_test

import (
	"log"
	"log/slog"
	"os"
	"testing"
	"wacdo-backend/config"
	"wacdo-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// load .env.test variables using go dot env
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}
	// InitDB handles errors itself, not in main.
	testDB = config.InitDB()
	if testDB == nil {
		slog.Error("Unable to instance DB object")
		os.Exit(1)
	}
	slog.Info("Successfully connected to", "database", os.Getenv("DB_NAME"))
	err = testDB.AutoMigrate(models.Menu{},
		models.Product{},
		models.MenuProduct{},
		models.Order{},
		models.OrderMenu{},
		models.User{},
		models.OrderProduct{},
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("Database setup is complete.")

	exitValue := m.Run() // This runs the tests and pass the test exist code to the variable
	// if code needs to run between the tests and the exit, it must go here
	testDB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	testDB.Exec("TRUNCATE TABLE order_products")
	testDB.Exec("TRUNCATE TABLE order_menus")
	testDB.Exec("TRUNCATE TABLE menus")
	testDB.Exec("TRUNCATE TABLE products")
	testDB.Exec("TRUNCATE TABLE menu_products")
	testDB.Exec("TRUNCATE TABLE users")
	testDB.Exec("TRUNCATE TABLE orders")
	testDB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	// Exit with the test return code
	os.Exit(exitValue)
}
