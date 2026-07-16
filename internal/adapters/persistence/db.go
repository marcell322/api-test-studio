package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewGormDB opens (and creates if necessary) a SQLite database at dbPath and returns a *gorm.DB.
// It ensures parent directories exist and returns clear wrapped errors.
func NewGormDB(dbPath string) (*gorm.DB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("empty db path")
	}
	// ensure parent dir exists
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create db directory %s: %w", dir, err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db %s: %w", dbPath, err)
	}
	return db, nil
}

// AutoMigrate runs GORM automigrations for the provided models. Returns descriptive errors.
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if len(models) == 0 {
		return fmt.Errorf("no models provided for automigrate")
	}
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("automigrate failed: %w", err)
	}
	return nil
}

// CloseDB closes the underlying sql.DB connection. Safe to call; returns wrapped errors.
func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("getting sql DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("closing db: %w", err)
	}
	return nil
}
