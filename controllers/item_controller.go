package controllers

import (
	"net/http"

	"inventory-api/models"
	"inventory-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemController struct {
	itemService *utils.ItemService
}

func NewItemController() *ItemController {
	return &ItemController{
		itemService: utils.NewItemService(),
	}
}

func (c *ItemController) SetItemService(service *utils.ItemService) {
	c.itemService = service
}

// CreateItem handles POST /inventory
// @Summary Create a new item
// @Description Create a new inventory item
// @Tags items
// @Accept json
// @Produce json
// @Param item body models.CreateItemRequest true "Item data"
// @Success 201 {object} models.Item
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory [post]
func (h *ItemController) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.CreateItem(&req)
	if err != nil {
		utils.Error.Printf("Failed to create item: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	utils.Info.Printf("Created item: %s", item.ID)
	c.JSON(http.StatusCreated, item)
}

// GetItem handles GET /inventory/:id
// @Summary Get an item by ID
// @Description Get a specific inventory item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.Item
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory/{id} [get]
func (h *ItemController) GetItem(c *gin.Context) {
	id := c.Param("id")
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		utils.Error.Printf("Invalid UUID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid UUID format",
			Message: "The provided ID is not a valid UUID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.GetItem(id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Item not found",
				Message: "The requested item does not exist",
				Code:    http.StatusNotFound,
			})
			return
		}
		
		utils.Error.Printf("Failed to get item: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// UpdateItem handles PUT /inventory/:id
// @Summary Update an item
// @Description Update an existing inventory item
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body models.UpdateItemRequest true "Updated item data"
// @Success 200 {object} models.Item
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory/{id} [put]
func (h *ItemController) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		utils.Error.Printf("Invalid UUID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid UUID format",
			Message: "The provided ID is not a valid UUID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	item, err := h.itemService.UpdateItem(id, &req)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Item not found",
				Message: "The requested item does not exist",
				Code:    http.StatusNotFound,
			})
			return
		}
		
		utils.Error.Printf("Failed to update item: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	utils.Info.Printf("Updated item: %s", item.ID)
	c.JSON(http.StatusOK, item)
}

// DeleteItem handles DELETE /inventory/:id
// @Summary Delete an item
// @Description Delete an inventory item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory/{id} [delete]
func (h *ItemController) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		utils.Error.Printf("Invalid UUID format: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid UUID format",
			Message: "The provided ID is not a valid UUID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.itemService.DeleteItem(id)
	if err != nil {
		if err.Error() == "item not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Item not found",
				Message: "The requested item does not exist",
				Code:    http.StatusNotFound,
			})
			return
		}
		
		utils.Error.Printf("Failed to delete item: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete item",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	utils.Info.Printf("Deleted item: %s", id)
	c.Status(http.StatusNoContent)
}

// GetItems handles GET /inventory
// @Summary Get all items
// @Description Get all inventory items with pagination, filtering, and sorting
// @Tags items
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (max 100)" default(10)
// @Param cursor query string false "Cursor for pagination"
// @Param name query string false "Filter by item name (partial match)"
// @Param min_stock query int false "Filter by minimum stock level"
// @Param min_price query number false "Filter by minimum price"
// @Param max_price query number false "Filter by maximum price"
// @Param sort_by query string false "Sort by field (name, stock, price, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory [get]
func (h *ItemController) GetItems(c *gin.Context) {
	// Parse pagination parameters
	var pagination models.PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.Error.Printf("Invalid pagination parameters: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid pagination parameters",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Parse filter parameters
	var filters models.FilterRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.Error.Printf("Invalid filter parameters: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid filter parameters",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Parse sort parameters
	var sort models.SortRequest
	if err := c.ShouldBindQuery(&sort); err != nil {
		utils.Error.Printf("Invalid sort parameters: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid sort parameters",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Set default values
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}
	if sort.SortBy == "" {
		sort.SortBy = "created_at"
	}
	if sort.SortOrder == "" {
		sort.SortOrder = "desc"
	}

	response, err := h.itemService.GetItems(&pagination, &filters, &sort)
	if err != nil {
		utils.Error.Printf("Failed to get items: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get items",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetItemStats handles GET /inventory/stats
// @Summary Get inventory statistics
// @Description Get statistics about the inventory
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory/stats [get]
func (h *ItemController) GetItemStats(c *gin.Context) {
	stats, err := h.itemService.GetItemStats()
	if err != nil {
		utils.Error.Printf("Failed to get item stats: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get item stats",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// SeedDatabase handles POST /inventory/seed
// @Summary Seed the database
// @Description Seed the database with sample data
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} models.ErrorResponse
// @Router /inventory/seed [post]
func (h *ItemController) SeedDatabase(c *gin.Context) {
	err := h.itemService.SeedDatabase()
	if err != nil {
		utils.Error.Printf("Failed to seed database: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to seed database",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	utils.Info.Println("Database seeded successfully")
	c.JSON(http.StatusOK, map[string]string{
		"message": "Database seeded successfully with sample data",
	})
}
