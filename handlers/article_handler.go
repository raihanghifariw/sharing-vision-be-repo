package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sharing-vision/sharing-vision-be/services"
	"github.com/sharing-vision/sharing-vision-be/validators"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	service services.ArticleService
}

func NewArticleHandler(service services.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req validators.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validators.ValidateArticle(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	limitStr := c.Param("limit")
	offsetStr := c.Param("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	// Optional ?status= query param. Validated if provided.
	status := c.Query("status")
	if status != "" {
		valid := map[string]bool{"publish": true, "draft": true, "thrash": true}
		if !valid[status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status must be publish, draft, or thrash"})
			return
		}
	}

	posts, err := h.service.FindAll(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	post, err := h.service.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	var req validators.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validators.ValidateArticle(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(uint(id), req); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
