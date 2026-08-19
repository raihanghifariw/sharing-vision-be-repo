package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase opens a GORM connection and returns an error instead of calling
// log.Fatalf — callers decide whether to retry or crash.
func InitDatabase() error {
	// Prefer direct Railway MySQL vars (MYSQLHOST, etc.) over DB_* aliases.
	// Falls back to DB_* so local .env still works.
	host := firstNonEmpty(os.Getenv("MYSQLHOST"), os.Getenv("DB_HOST"))
	port := firstNonEmpty(os.Getenv("MYSQLPORT"), os.Getenv("DB_PORT"), "3306")
	user := firstNonEmpty(os.Getenv("MYSQLUSER"), os.Getenv("DB_USER"))
	pass := firstNonEmpty(os.Getenv("MYSQLPASSWORD"), os.Getenv("DB_PASSWORD"))
	name := firstNonEmpty(os.Getenv("MYSQLDATABASE"), os.Getenv("DB_NAME"))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}

	DB = db
	return nil
}

// firstNonEmpty returns the first non-empty string from the arguments.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
