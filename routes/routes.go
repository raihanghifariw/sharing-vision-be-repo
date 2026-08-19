package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sharing-vision/sharing-vision-be/handlers"
)

// SetupRoutes registers CORS, security headers, and the health check.
// Article routes are registered separately via RegisterArticleRoutes after DB is ready.
func SetupRoutes(router *gin.Engine, articleHandler *handlers.ArticleHandler) {
	applyMiddleware(router)
	RegisterArticleRoutes(router, articleHandler)
}

// RegisterArticleRoutes adds the five article endpoints.
// Called from main.go once the DB connection is confirmed.
func RegisterArticleRoutes(router *gin.Engine, articleHandler *handlers.ArticleHandler) {
	router.POST("/article/", articleHandler.CreateArticle)
	router.GET("/article/:limit/:offset", articleHandler.GetArticles)
	router.GET("/article/detail/:id", articleHandler.GetArticleByID)
	router.PUT("/article/:id", articleHandler.UpdateArticle)
	router.DELETE("/article/:id", articleHandler.DeleteArticle)
}

func applyMiddleware(router *gin.Engine) {
	// Security headers
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	})

	// CORS — allow all origins so any frontend deployment can connect.
	// This is safe for a public read/write article API with no auth.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))
}
