package routes

import (
	"log"

	"github.com/esrid/watcher"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterWatcher(db *gorm.DB, router *gin.Engine) {
	// Register the schema introspection watcher
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	inspector, err := watcher.NewInspector(sqlDB)
	if err != nil {
		log.Fatal(err)
	}
	// Register the watcher at /schema endpoint
	routes := router.Group("/schema")
	// use gin.WrapF to wrap the http handler as a gin.HandlerFunc and serve it through Gin.
	routes.GET("", gin.WrapF(watcher.HTTPHandler(inspector)))

}
