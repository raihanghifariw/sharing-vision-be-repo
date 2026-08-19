package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sharing-vision/sharing-vision-be/config"
	"github.com/sharing-vision/sharing-vision-be/handlers"
	"github.com/sharing-vision/sharing-vision-be/migrations"
	"github.com/sharing-vision/sharing-vision-be/repositories"
	"github.com/sharing-vision/sharing-vision-be/routes"
	"github.com/sharing-vision/sharing-vision-be/services"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// Init DB
	config.InitDatabase()

	// Run migrations
	migrations.RunMigrations(config.DB)

	// Wire dependencies
	articleRepo := repositories.NewArticleRepository(config.DB)
	articleService := services.NewArticleService(articleRepo)
	articleHandler := handlers.NewArticleHandler(articleService)

	// Setup router
	router := gin.Default()
	routes.SetupRoutes(router, articleHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
