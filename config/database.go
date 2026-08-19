package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase opens a GORM connection and returns an error instead of calling
// log.Fatalf — callers decide whether to retry or crash.
//
// Priority order for connection info:
//  1. MYSQL_URL or DATABASE_URL  (Railway auto-injects this — most reliable)
//  2. Individual MYSQLHOST/MYSQLPORT/… vars
//  3. DB_HOST/DB_PORT/… vars (local .env fallback)
func InitDatabase() error {
	dsn, err := buildDSN()
	if err != nil {
		return err
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open: %w", err)
	}

	DB = db
	return nil
}

func buildDSN() (string, error) {
	// 1. Try MYSQL_URL / DATABASE_URL first (Railway injects these automatically)
	for _, key := range []string{"MYSQL_URL", "DATABASE_URL"} {
		if raw := os.Getenv(key); raw != "" {
			dsn, err := urlToDSN(raw)
			if err == nil {
				return dsn, nil
			}
		}
	}

	// 2. Individual env vars — Railway MYSQL* or local DB_*
	host := firstNonEmpty(os.Getenv("MYSQLHOST"), os.Getenv("DB_HOST"))
	port := firstNonEmpty(os.Getenv("MYSQLPORT"), os.Getenv("DB_PORT"), "3306")
	user := firstNonEmpty(os.Getenv("MYSQLUSER"), os.Getenv("DB_USER"))
	pass := firstNonEmpty(os.Getenv("MYSQLPASSWORD"), os.Getenv("DB_PASSWORD"))
	name := firstNonEmpty(os.Getenv("MYSQLDATABASE"), os.Getenv("DB_NAME"))

	if host == "" || user == "" || name == "" {
		return "", fmt.Errorf("database config missing: host=%q user=%q name=%q", host, user, name)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name,
	), nil
}

// urlToDSN converts a mysql://user:pass@host:port/dbname URL to a GORM DSN.
func urlToDSN(raw string) (string, error) {
	// Railway sometimes uses "mysql://" scheme
	if !strings.HasPrefix(raw, "mysql://") && !strings.HasPrefix(raw, "mysql2://") {
		return "", fmt.Errorf("not a mysql URL: %s", raw)
	}
	raw = strings.Replace(raw, "mysql2://", "mysql://", 1)

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	pass, _ := u.User.Password()
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "3306"
	}
	name := strings.TrimPrefix(u.Path, "/")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		u.User.Username(), pass, host, port, name,
	), nil
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
