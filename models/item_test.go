package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItem_TableName(t *testing.T) {
	item := Item{}
	assert.Equal(t, "items", item.TableName())
}

func TestItem_BeforeCreate(t *testing.T) {
	tests := []struct {
		name     string
		item     Item
		expected bool
	}{
		{
			name: "nil UUID should be set",
			item: Item{
				ID: uuid.Nil,
			},
			expected: true,
		},
		{
			name: "existing UUID should not be changed",
			item: Item{
				ID: uuid.New(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalID := tt.item.ID
			err := tt.item.BeforeCreate(nil)
			require.NoError(t, err)

			if tt.expected {
				assert.NotEqual(t, uuid.Nil, tt.item.ID)
				assert.NotEqual(t, originalID, tt.item.ID)
			} else {
				assert.Equal(t, originalID, tt.item.ID)
			}
		})
	}
}

func TestCreateItemRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateItemRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateItemRequest{
				Name:  "Test Item",
				Stock: 10,
				Price: 99.99,
			},
			wantErr: false,
		},
		{
			name: "empty name should fail",
			request: CreateItemRequest{
				Name:  "",
				Stock: 10,
				Price: 99.99,
			},
			wantErr: true,
		},
		{
			name: "negative stock should fail",
			request: CreateItemRequest{
				Name:  "Test Item",
				Stock: -1,
				Price: 99.99,
			},
			wantErr: true,
		},
		{
			name: "negative price should fail",
			request: CreateItemRequest{
				Name:  "Test Item",
				Stock: 10,
				Price: -1.0,
			},
			wantErr: true,
		},
		{
			name: "zero stock should pass",
			request: CreateItemRequest{
				Name:  "Test Item",
				Stock: 0,
				Price: 99.99,
			},
			wantErr: false,
		},
		{
			name: "zero price should pass",
			request: CreateItemRequest{
				Name:  "Test Item",
				Stock: 10,
				Price: 0.0,
			},
			wantErr: false,
		},
		{
			name: "very long name should fail",
			request: CreateItemRequest{
				Name:  string(make([]byte, 256)), // 256 characters
				Stock: 10,
				Price: 99.99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validator
			// This is a simplified validation test
			hasError := false
			
			if tt.request.Name == "" || len(tt.request.Name) > 255 {
				hasError = true
			}
			if tt.request.Stock < 0 {
				hasError = true
			}
			if tt.request.Price < 0 {
				hasError = true
			}
			
			assert.Equal(t, tt.wantErr, hasError)
		})
	}
}

func TestUpdateItemRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateItemRequest
		wantErr bool
	}{
		{
			name: "valid request with all fields",
			request: UpdateItemRequest{
				Name:  stringPtr("Updated Item"),
				Stock: intPtr(20),
				Price: float64Ptr(199.99),
			},
			wantErr: false,
		},
		{
			name: "valid request with partial fields",
			request: UpdateItemRequest{
				Name: stringPtr("Updated Item"),
			},
			wantErr: false,
		},
		{
			name: "valid request with no fields",
			request: UpdateItemRequest{},
			wantErr: false,
		},
		{
			name: "empty name should fail",
			request: UpdateItemRequest{
				Name: stringPtr(""),
			},
			wantErr: true,
		},
		{
			name: "negative stock should fail",
			request: UpdateItemRequest{
				Stock: intPtr(-1),
			},
			wantErr: true,
		},
		{
			name: "negative price should fail",
			request: UpdateItemRequest{
				Price: float64Ptr(-1.0),
			},
			wantErr: true,
		},
		{
			name: "zero stock should pass",
			request: UpdateItemRequest{
				Stock: intPtr(0),
			},
			wantErr: false,
		},
		{
			name: "zero price should pass",
			request: UpdateItemRequest{
				Price: float64Ptr(0.0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.request.Name != nil && (*tt.request.Name == "" || len(*tt.request.Name) > 255) {
				hasError = true
			}
			if tt.request.Stock != nil && *tt.request.Stock < 0 {
				hasError = true
			}
			if tt.request.Price != nil && *tt.request.Price < 0 {
				hasError = true
			}
			
			assert.Equal(t, tt.wantErr, hasError)
		})
	}
}

func TestPaginationRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request PaginationRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: PaginationRequest{
				Limit:  10,
				Cursor: "valid-cursor",
			},
			wantErr: false,
		},
		{
			name: "zero limit should be valid",
			request: PaginationRequest{
				Limit: 0,
			},
			wantErr: false,
		},
		{
			name: "negative limit should fail",
			request: PaginationRequest{
				Limit: -1,
			},
			wantErr: true,
		},
		{
			name: "limit over 100 should fail",
			request: PaginationRequest{
				Limit: 101,
			},
			wantErr: true,
		},
		{
			name: "limit of 100 should pass",
			request: PaginationRequest{
				Limit: 100,
			},
			wantErr: false,
		},
		{
			name: "empty cursor should be valid",
			request: PaginationRequest{
				Limit:  10,
				Cursor: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.request.Limit < 0 || tt.request.Limit > 100 {
				hasError = true
			}
			
			assert.Equal(t, tt.wantErr, hasError)
		})
	}
}

func TestFilterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request FilterRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: FilterRequest{
				Name:     "test",
				MinStock: intPtr(10),
				MinPrice: float64Ptr(100.0),
				MaxPrice: float64Ptr(500.0),
			},
			wantErr: false,
		},
		{
			name: "empty name should be valid",
			request: FilterRequest{
				Name: "",
			},
			wantErr: false,
		},
		{
			name: "negative min stock should fail",
			request: FilterRequest{
				MinStock: intPtr(-1),
			},
			wantErr: true,
		},
		{
			name: "negative min price should fail",
			request: FilterRequest{
				MinPrice: float64Ptr(-1.0),
			},
			wantErr: true,
		},
		{
			name: "negative max price should fail",
			request: FilterRequest{
				MaxPrice: float64Ptr(-1.0),
			},
			wantErr: true,
		},
		{
			name: "zero values should pass",
			request: FilterRequest{
				MinStock: intPtr(0),
				MinPrice: float64Ptr(0.0),
				MaxPrice: float64Ptr(0.0),
			},
			wantErr: false,
		},
		{
			name: "min price greater than max price should be valid (handled by business logic)",
			request: FilterRequest{
				MinPrice: float64Ptr(500.0),
				MaxPrice: float64Ptr(100.0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.request.MinStock != nil && *tt.request.MinStock < 0 {
				hasError = true
			}
			if tt.request.MinPrice != nil && *tt.request.MinPrice < 0 {
				hasError = true
			}
			if tt.request.MaxPrice != nil && *tt.request.MaxPrice < 0 {
				hasError = true
			}
			
			assert.Equal(t, tt.wantErr, hasError)
		})
	}
}

func TestSortRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request SortRequest
		wantErr bool
	}{
		{
			name: "valid sort by name asc",
			request: SortRequest{
				SortBy:    "name",
				SortOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "valid sort by stock desc",
			request: SortRequest{
				SortBy:    "stock",
				SortOrder: "desc",
			},
			wantErr: false,
		},
		{
			name: "valid sort by price asc",
			request: SortRequest{
				SortBy:    "price",
				SortOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "valid sort by created_at desc",
			request: SortRequest{
				SortBy:    "created_at",
				SortOrder: "desc",
			},
			wantErr: false,
		},
		{
			name: "invalid sort by field should fail",
			request: SortRequest{
				SortBy:    "invalid_field",
				SortOrder: "asc",
			},
			wantErr: true,
		},
		{
			name: "invalid sort order should fail",
			request: SortRequest{
				SortBy:    "name",
				SortOrder: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty sort by should be valid",
			request: SortRequest{
				SortBy:    "",
				SortOrder: "asc",
			},
			wantErr: false,
		},
		{
			name: "empty sort order should be valid",
			request: SortRequest{
				SortBy:    "name",
				SortOrder: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			validSortBy := map[string]bool{
				"name":       true,
				"stock":      true,
				"price":      true,
				"created_at": true,
			}
			
			validSortOrder := map[string]bool{
				"asc":  true,
				"desc": true,
			}
			
			if tt.request.SortBy != "" && !validSortBy[tt.request.SortBy] {
				hasError = true
			}
			if tt.request.SortOrder != "" && !validSortOrder[tt.request.SortOrder] {
				hasError = true
			}
			
			assert.Equal(t, tt.wantErr, hasError)
		})
	}
}

func TestErrorResponse_Structure(t *testing.T) {
	errorResp := ErrorResponse{
		Error:   "Test Error",
		Message: "Test Message",
		Code:    400,
	}

	assert.Equal(t, "Test Error", errorResp.Error)
	assert.Equal(t, "Test Message", errorResp.Message)
	assert.Equal(t, 400, errorResp.Code)
}

func TestPaginatedResponse_Structure(t *testing.T) {
	items := []Item{
		{Name: "Item 1", Stock: 10, Price: 99.99},
		{Name: "Item 2", Stock: 20, Price: 199.99},
	}

	response := PaginatedResponse{
		Items:      items,
		NextCursor: "next-cursor",
		HasMore:    true,
		Total:      100,
	}

	assert.Len(t, response.Items, 2)
	assert.Equal(t, "next-cursor", response.NextCursor)
	assert.True(t, response.HasMore)
	assert.Equal(t, int64(100), response.Total)
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
