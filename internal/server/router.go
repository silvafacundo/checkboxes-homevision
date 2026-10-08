// Package server wires up the HTTP router.
package server

import (
	"github.com/gin-gonic/gin"

	"homevision/internal/config"
	"homevision/internal/detect"
	"homevision/internal/health"
)

// NewRouter builds the Gin engine with all routes registered.
func NewRouter(cfg config.Config) *gin.Engine {
	router := gin.Default()

	// Add simple CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve the local index.html file at the root URL
	router.StaticFile("/", "./public/index.html")

	// Group all API routes under /api
	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Checkbox detector"})
		})
		api.GET("/ping", health.Ping)

		detectHandler := detect.NewDetectHandler(detect.NewDetector(), cfg.MaxUploadSize, cfg.AllowQueryOptions)
		api.POST("/detect", detectHandler.Detect)
	}

	return router
}
