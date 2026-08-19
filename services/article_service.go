package services

import (
	"github.com/sharing-vision/sharing-vision-be/models"
	"github.com/sharing-vision/sharing-vision-be/repositories"
	"github.com/sharing-vision/sharing-vision-be/validators"
)

type ArticleService interface {
	Create(req validators.ArticleRequest) error
	FindAll(limit, offset int) ([]models.Post, error)
	FindByID(id uint) (*models.Post, error)
	Update(id uint, req validators.ArticleRequest) error
	Delete(id uint) error
}

type articleService struct {
	repo repositories.ArticleRepository
}

func NewArticleService(repo repositories.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) Create(req validators.ArticleRequest) error {
	post := &models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}
	return s.repo.Create(post)
}

func (s *articleService) FindAll(limit, offset int) ([]models.Post, error) {
	return s.repo.FindAll(limit, offset)
}

func (s *articleService) FindByID(id uint) (*models.Post, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) Update(id uint, req validators.ArticleRequest) error {
	post := &models.Post{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}
	return s.repo.Update(id, post)
}

func (s *articleService) Delete(id uint) error {
	return s.repo.Delete(id)
}
