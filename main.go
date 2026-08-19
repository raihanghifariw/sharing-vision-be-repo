package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	// Load .env (ignored on Railway — env vars are injected directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server immediately so Railway healthcheck can reach /health
	// DB connection happens in a goroutine with retry; /health returns 200 right away.
	router := gin.Default()

	// Health check — always responds, even before DB is ready
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Connect DB with retry in background, then register article routes
	go func() {
		for attempt := 1; attempt <= 10; attempt++ {
			log.Printf("DB connect attempt %d/10...", attempt)
			if err := config.InitDatabase(); err != nil {
				log.Printf("DB connect failed: %v — retrying in 5s", err)
				time.Sleep(5 * time.Second)
				continue
			}

			migrations.RunMigrations(config.DB)

			articleRepo := repositories.NewArticleRepository(config.DB)
			articleService := services.NewArticleService(articleRepo)
			articleHandler := handlers.NewArticleHandler(articleService)

			routes.RegisterArticleRoutes(router, articleHandler)
			log.Println("Article routes registered — backend fully ready")
			return
		}
		log.Println("ERROR: Could not connect to database after 10 attempts")
	}()

	log.Printf("Server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
