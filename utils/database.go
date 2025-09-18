package utils

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *Config) error {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}

// Migrate runs database migrations (development mode only)
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Run migration files in order
	migrationFiles := []string{
		"migrations/001_drop_tables.sql",
		"migrations/002_create_items_table.sql",
	}

	for _, file := range migrationFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if err := DB.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		Info.Printf("Successfully executed migration: %s", file)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Health checks the database connection health
func Health() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
