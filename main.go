package main

import (
	"log"
	"os"
	"rwa-backend/database"
	"rwa-backend/handlers"
	"rwa-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// Create Gin router
	r := gin.Default()

	// Setup CORS middleware
	r.Use(middleware.SetupCORS())

	// Setup routes
	setupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 RWA Backend starting on port %s", port)
	log.Printf("📊 Health check: http://localhost:%s/api/health", port)
	log.Printf("📋 Assets API: http://localhost:%s/api/assets", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// setupRoutes configures all API routes
func setupRoutes(r *gin.Engine) {
	// Basic health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "RWA Backend API is running! 🚀",
			"version": "1.0.0",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", handlers.HealthCheck)

		// Asset routes
		assets := api.Group("/assets")
		{
			assets.GET("/", handlers.GetAllAssets)                        // GET /api/assets
			assets.POST("/mint", handlers.MintAsset)                      // POST /api/assets/mint
			assets.GET("/:id", handlers.GetAssetByID)                     // GET /api/assets/:id
			assets.PATCH("/:id/contract", handlers.UpdateContractAddress) // PATCH /api/assets/:id/contract
			assets.GET("/:id/stats", handlers.GetAssetStats)              // GET /api/assets/:id/stats
			assets.PATCH("/:id/staking", handlers.UpdateAssetStaking)     // PATCH /api/assets/:id/staking (for testing)
		}
	}
}
