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
	// Load .env for local dev (ignored on Railway — env vars are injected directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Log DB config on startup so we can diagnose missing vars immediately
	log.Printf("DB config: MYSQL_URL=%q MYSQLHOST=%q MYSQLPORT=%q MYSQLUSER=%q MYSQLDATABASE=%q",
		maskSecret(os.Getenv("MYSQL_URL")),
		os.Getenv("MYSQLHOST"),
		os.Getenv("MYSQLPORT"),
		os.Getenv("MYSQLUSER"),
		os.Getenv("MYSQLDATABASE"),
	)

	// Start HTTP server immediately so Railway healthcheck passes right away.
	// DB connects in background goroutine; article routes register once DB is ready.
	router := gin.Default()

	// Health check — always responds, even before DB is ready
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Apply CORS + security headers now so they work for all routes including article routes
	routes.SetupMiddleware(router)

	// Connect DB with retry in background, then register article routes
	go func() {
		for attempt := 1; attempt <= 20; attempt++ {
			log.Printf("DB connect attempt %d/20...", attempt)
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
		log.Println("FATAL: Could not connect to database after 20 attempts — article routes NOT registered")
	}()

	log.Printf("Server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// maskSecret shows only the scheme part of a URL to avoid leaking credentials in logs.
func maskSecret(s string) string {
	if s == "" {
		return ""
	}
	if len(s) > 12 {
		return s[:12] + "***"
	}
	return "***"
}
