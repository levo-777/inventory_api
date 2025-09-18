package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key" swaggertype:"string" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string         `json:"name" gorm:"not null;size:255" binding:"required,min=1,max=255" example:"Laptop"`
	Stock     int            `json:"stock" gorm:"not null;default:0" binding:"required,min=0" example:"50"`
	Price     float64        `json:"price" gorm:"not null;type:decimal(10,2)" binding:"required,min=0" example:"999.99"`
	CreatedAt time.Time      `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time      `json:"updated_at" swaggertype:"string" format:"date-time"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time"`
}

// TableName returns the table name for the Item model
func (Item) TableName() string {
	return "items"
}

// BeforeCreate hook to generate UUID if not set
func (i *Item) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

// CreateItemRequest represents the request payload for creating an item
type CreateItemRequest struct {
	Name  string  `json:"name" binding:"required,min=1,max=255" example:"Laptop"`
	Stock int     `json:"stock" binding:"required,min=0" example:"50"`
	Price float64 `json:"price" binding:"required,min=0" example:"999.99"`
}

// UpdateItemRequest represents the request payload for updating an item
type UpdateItemRequest struct {
	Name  *string  `json:"name,omitempty" binding:"omitempty,min=1,max=255" example:"Updated Laptop"`
	Stock *int     `json:"stock,omitempty" binding:"omitempty,min=0" example:"75"`
	Price *float64 `json:"price,omitempty" binding:"omitempty,min=0" example:"1099.99"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100" example:"10"`
	Cursor string `form:"cursor" example:"eyJpZCI6IjU1MGU4NDAwLWUyOWItNDFkNC1hNzE2LTQ0NjY1NTQ0MDAwMCJ9"`
}

// FilterRequest represents filtering parameters
type FilterRequest struct {
	Name      string `form:"name" example:"laptop"`
	MinStock  *int   `form:"min_stock" binding:"omitempty,min=0" example:"10"`
	MinPrice  *float64 `form:"min_price" binding:"omitempty,min=0" example:"100.0"`
	MaxPrice  *float64 `form:"max_price" binding:"omitempty,min=0" example:"2000.0"`
}

// SortRequest represents sorting parameters
type SortRequest struct {
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name stock price created_at" example:"name"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc" example:"asc"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Items      []Item `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
	Total      int64  `json:"total,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
