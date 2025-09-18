package integrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"inventory-api/controllers"
	"inventory-api/models"
	"inventory-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemHandler_CreateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	// Create handler with test database
	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.POST("/inventory", handler.CreateItem)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid item creation",
			requestBody: models.CreateItemRequest{
				Name:  "Test Laptop",
				Stock: 50,
				Price: 999.99,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "item with zero stock",
			requestBody: models.CreateItemRequest{
				Name:  "Test Item",
				Stock: 0,
				Price: 99.99,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "item with zero price",
			requestBody: models.CreateItemRequest{
				Name:  "Free Item",
				Stock: 10,
				Price: 0.0,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"name": "Test Item",
				// Missing stock and price
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "negative stock",
			requestBody: models.CreateItemRequest{
				Name:  "Test Item",
				Stock: -1,
				Price: 99.99,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "negative price",
			requestBody: models.CreateItemRequest{
				Name:  "Test Item",
				Stock: 10,
				Price: -1.0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "empty name",
			requestBody: models.CreateItemRequest{
				Name:  "",
				Stock: 10,
				Price: 99.99,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				err = json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
			} else {
				var item models.Item
				err = json.Unmarshal(w.Body.Bytes(), &item)
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, item.ID)
			}
		})
	}
}

func TestItemHandler_GetItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.GET("/inventory/:id", handler.GetItem)

	// Create a test item
	item := testDB.CreateTestItem(t, "Test Item", 10, 99.99)

	tests := []struct {
		name           string
		itemID         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "existing item",
			itemID:         item.ID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent item",
			itemID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
		{
			name:           "invalid UUID",
			itemID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid UUID format",
		},
		{
			name:           "empty ID",
			itemID:         "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/inventory/"+tt.itemID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
			} else {
				var responseItem models.Item
				err := json.Unmarshal(w.Body.Bytes(), &responseItem)
				require.NoError(t, err)
				assert.Equal(t, item.ID, responseItem.ID)
			}
		})
	}
}

func TestItemHandler_UpdateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.PUT("/inventory/:id", handler.UpdateItem)

	// Create a test item
	item := testDB.CreateTestItem(t, "Original Item", 10, 99.99)

	tests := []struct {
		name           string
		itemID         string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "update name only",
			itemID: item.ID.String(),
			requestBody: models.UpdateItemRequest{
				Name: utils.StringPtr("Updated Item"),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "update all fields",
			itemID: item.ID.String(),
			requestBody: models.UpdateItemRequest{
				Name:  utils.StringPtr("Fully Updated Item"),
				Stock: utils.IntPtr(20),
				Price: utils.Float64Ptr(199.99),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "update non-existent item",
			itemID:         uuid.New().String(),
			requestBody:    models.UpdateItemRequest{Name: utils.StringPtr("Updated Item")},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
		{
			name:           "invalid UUID",
			itemID:         "invalid-uuid",
			requestBody:    models.UpdateItemRequest{Name: utils.StringPtr("Updated Item")},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid UUID format",
		},
		{
			name:           "invalid JSON",
			itemID:         item.ID.String(),
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:   "empty name",
			itemID: item.ID.String(),
			requestBody: models.UpdateItemRequest{
				Name: utils.StringPtr(""),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name:   "negative stock",
			itemID: item.ID.String(),
			requestBody: models.UpdateItemRequest{
				Stock: utils.IntPtr(-1),
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPut, "/inventory/"+tt.itemID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				err = json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
			} else {
				var responseItem models.Item
				err = json.Unmarshal(w.Body.Bytes(), &responseItem)
				require.NoError(t, err)
				assert.Equal(t, item.ID, responseItem.ID)
			}
		})
	}
}

func TestItemHandler_DeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.DELETE("/inventory/:id", handler.DeleteItem)

	// Create a test item
	item := testDB.CreateTestItem(t, "Test Item", 10, 99.99)

	tests := []struct {
		name           string
		itemID         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "delete existing item",
			itemID:         item.ID.String(),
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existent item",
			itemID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			expectedError:  "Item not found",
		},
		{
			name:           "invalid UUID",
			itemID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/inventory/"+tt.itemID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
			} else {
				// For successful deletion, body should be empty
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func TestItemHandler_GetItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.GET("/inventory", handler.GetItems)

	// Create test items
	testDB.CreateTestItem(t, "Laptop Pro", 50, 999.99)
	testDB.CreateTestItem(t, "Gaming Mouse", 100, 49.99)
	testDB.CreateTestItem(t, "Wireless Keyboard", 25, 79.99)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
		expectedError  string
	}{
		{
			name:           "get all items",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "get items with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "get items with name filter",
			queryParams:    "?name=Laptop",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "get items with min stock filter",
			queryParams:    "?min_stock=50",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "get items with price range filter",
			queryParams:    "?min_price=50&max_price=100",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "get items with sorting",
			queryParams:    "?sort_by=name&sort_order=asc",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "get items with combined filters",
			queryParams:    "?name=Mouse&min_stock=50&sort_by=price&sort_order=desc",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "invalid limit",
			queryParams:    "?limit=-1",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid pagination parameters",
		},
		{
			name:           "limit too high",
			queryParams:    "?limit=101",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid pagination parameters",
		},
		{
			name:           "invalid sort field",
			queryParams:    "?sort_by=invalid_field",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid sort parameters",
		},
		{
			name:           "invalid sort order",
			queryParams:    "?sort_order=invalid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid sort parameters",
		},
		{
			name:           "invalid min stock",
			queryParams:    "?min_stock=-1",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid filter parameters",
		},
		{
			name:           "invalid min price",
			queryParams:    "?min_price=-1",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid filter parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/inventory"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var errorResp models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, errorResp.Error)
			} else {
				var response models.PaginatedResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Len(t, response.Items, tt.expectedCount)
			}
		})
	}
}

func TestItemHandler_GetItemStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.GET("/inventory/stats", handler.GetItemStats)

	// Create test items
	testDB.CreateTestItem(t, "Item 1", 10, 100.0)
	testDB.CreateTestItem(t, "Item 2", 5, 200.0) // Low stock
	testDB.CreateTestItem(t, "Item 3", 15, 300.0)

	req := httptest.NewRequest(http.MethodGet, "/inventory/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var stats map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &stats)
	require.NoError(t, err)

	assert.Equal(t, float64(3), stats["total_items"])
	assert.Equal(t, float64(200), stats["average_price"])
	assert.Equal(t, float64(1), stats["low_stock_items"])
	assert.Equal(t, float64(6000), stats["total_value"])
}

func TestItemHandler_SeedDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.POST("/inventory/seed", handler.SeedDatabase)

	req := httptest.NewRequest(http.MethodPost, "/inventory/seed", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Database seeded successfully with sample data", response["message"])
}

func TestItemHandler_ConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := utils.SetupTestRouter()

	testDB := utils.NewTestDB(t)
	defer testDB.Close()

	handler := controllers.NewItemController()
	handler.SetItemService(utils.NewItemServiceWithDB(testDB.DB))

	router.POST("/inventory", handler.CreateItem)
	router.GET("/inventory", handler.GetItems)

	// Test concurrent creates
	const numRequests = 10
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			reqBody := models.CreateItemRequest{
				Name:  fmt.Sprintf("Concurrent Item %d", i),
				Stock: i + 1,
				Price: float64(i+1) * 10.0,
			}

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				errors <- fmt.Errorf("expected status 201, got %d", w.Code)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		assert.NoError(t, err)
	}

	// Verify all items were created
	req := httptest.NewRequest(http.MethodGet, "/inventory", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PaginatedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Items, numRequests)
}
