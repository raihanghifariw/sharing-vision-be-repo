package repositories

import (
	"github.com/sharing-vision/sharing-vision-be/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(post *models.Post) error
	// FindAll fetches posts. Pass status="" to return all statuses.
	FindAll(limit, offset int, status string) ([]models.Post, error)
	FindByID(id uint) (*models.Post, error)
	Update(id uint, post *models.Post) error
	Delete(id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *articleRepository) FindAll(limit, offset int, status string) ([]models.Post, error) {
	var posts []models.Post
	q := r.db.Limit(limit).Offset(offset)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&posts).Error
	return posts, err
}

func (r *articleRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *articleRepository) Update(id uint, post *models.Post) error {
	result := r.db.Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"title":    post.Title,
		"content":  post.Content,
		"category": post.Category,
		"status":   post.Status,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *articleRepository) Delete(id uint) error {
	result := r.db.Where("id = ?", id).Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
