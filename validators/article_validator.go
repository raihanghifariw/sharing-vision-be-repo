package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ArticleRequest struct {
	Title    string `json:"title"    validate:"required,min=20"`
	Content  string `json:"content"  validate:"required,min=200"`
	Category string `json:"category" validate:"required,min=3"`
	Status   string `json:"status"   validate:"required,oneof=publish draft thrash"`
}

var validate = validator.New()

func ValidateArticle(req ArticleRequest) error {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.Tag()
		param := e.Param()

		switch field {
		case "Title":
			if tag == "required" {
				return fmt.Errorf("title is required")
			}
			return fmt.Errorf("title minimal %s karakter", param)
		case "Content":
			if tag == "required" {
				return fmt.Errorf("content is required")
			}
			return fmt.Errorf("content minimal %s karakter", param)
		case "Category":
			if tag == "required" {
				return fmt.Errorf("category is required")
			}
			return fmt.Errorf("category minimal %s karakter", param)
		case "Status":
			if tag == "required" {
				return fmt.Errorf("status is required")
			}
			return fmt.Errorf("status must be publish, draft, or thrash")
		}
	}

	return fmt.Errorf("validation error")
}
