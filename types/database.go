package types

import (
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB(dataDir string) (*gorm.DB, error) {
	dbFile := filepath.Join(dataDir, "feathr.db")
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Automatically create schema if it doesn't exist
	db.AutoMigrate()

	// Close database connection
	defer db.Close()
	return db, nil
}
