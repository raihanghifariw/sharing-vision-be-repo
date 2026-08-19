package migrations

import (
	"log"

	"github.com/sharing-vision/sharing-vision-be/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed: table 'posts' ready")
}
