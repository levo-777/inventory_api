package utils

import (
	"testing"

	"inventory-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDB wraps a test database connection
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates a new in-memory SQLite database for testing
func NewTestDB(t *testing.T) *TestDB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.Item{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return &TestDB{DB: db}
}

// Close closes the test database connection
func (tdb *TestDB) Close() {
	if db, err := tdb.DB.DB(); err == nil {
		db.Close()
	}
}

// CreateTestItem creates a test item in the database
func (tdb *TestDB) CreateTestItem(t *testing.T, name string, stock int, price float64) *models.Item {
	item := &models.Item{
		ID:    uuid.New(),
		Name:  name,
		Stock: stock,
		Price: price,
	}

	if err := tdb.DB.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	return item
}

// SetupTestRouter creates a new Gin router for testing
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}
