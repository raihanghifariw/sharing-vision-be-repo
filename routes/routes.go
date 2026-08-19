package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sharing-vision/sharing-vision-be/handlers"
)

func SetupRoutes(router *gin.Engine, articleHandler *handlers.ArticleHandler) {
	// Security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Article routes
	router.POST("/article/", articleHandler.CreateArticle)
	router.GET("/article/:limit/:offset", articleHandler.GetArticles)
	router.GET("/article/detail/:id", articleHandler.GetArticleByID)
	router.PUT("/article/:id", articleHandler.UpdateArticle)
	router.DELETE("/article/:id", articleHandler.DeleteArticle)
}
