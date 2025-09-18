package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"inventory-api/models"

	"github.com/dgraph-io/ristretto/v2"
	"gorm.io/gorm"
)

type ItemService struct {
	db    *gorm.DB
	cache *ristretto.Cache[string, *models.Item]
}

type CursorData struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
}

func NewItemService() *ItemService {
	cache, err := ristretto.NewCache(&ristretto.Config[string, *models.Item]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		Error.Printf("Failed to create cache: %v", err)
		return &ItemService{db: DB}
	}

	return &ItemService{
		db:    DB,
		cache: cache,
	}
}

func NewItemServiceWithDB(db *gorm.DB) *ItemService {
	cache, err := ristretto.NewCache(&ristretto.Config[string, *models.Item]{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		Error.Printf("Failed to create cache: %v", err)
		return &ItemService{db: db}
	}

	return &ItemService{
		db:    db,
		cache: cache,
	}
}

func (s *ItemService) CreateItem(req *models.CreateItemRequest) (*models.Item, error) {
	item := &models.Item{
		Name:  req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	s.invalidateCache()

	return item, nil
}

func (s *ItemService) GetItem(id string) (*models.Item, error) {
	if item := s.getFromCache(id); item != nil {
		return item, nil
	}

	item := &models.Item{}
	if err := s.db.Where("id = ?", id).First(item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	s.setCache(id, item)

	return item, nil
}

func (s *ItemService) UpdateItem(id string, req *models.UpdateItemRequest) (*models.Item, error) {
	item := &models.Item{}
	if err := s.db.Where("id = ?", id).First(item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Stock != nil {
		item.Stock = *req.Stock
	}
	if req.Price != nil {
		item.Price = *req.Price
	}

	if err := s.db.Save(item).Error; err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	s.invalidateCache()

	return item, nil
}

func (s *ItemService) DeleteItem(id string) error {
	result := s.db.Where("id = ?", id).Delete(&models.Item{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	s.invalidateCache()

	return nil
}

func (s *ItemService) GetItems(pagination *models.PaginationRequest, filters *models.FilterRequest, sort *models.SortRequest) (*models.PaginatedResponse, error) {
	query := s.db.Model(&models.Item{})

	if filters != nil {
		if filters.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
		}
		if filters.MinStock != nil {
			query = query.Where("stock >= ?", *filters.MinStock)
		}
		if filters.MinPrice != nil {
			query = query.Where("price >= ?", *filters.MinPrice)
		}
		if filters.MaxPrice != nil {
			query = query.Where("price <= ?", *filters.MaxPrice)
		}
	}

	if sort != nil && sort.SortBy != "" {
		order := "ASC"
		if sort.SortOrder == "desc" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sort.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count items: %w", err)
	}

	var items []models.Item
	var nextCursor string
	var hasMore bool

	if pagination != nil && pagination.Cursor != "" {
		cursorData, err := s.decodeCursor(pagination.Cursor)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}

		query = query.Where("(created_at < ?) OR (created_at = ? AND id < ?)",
			cursorData.CreatedAt, cursorData.CreatedAt, cursorData.ID)
	}

	limit := 10
	if pagination != nil && pagination.Limit > 0 {
		limit = pagination.Limit
	}
	query = query.Limit(limit + 1)

	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	if len(items) > limit {
		hasMore = true
		items = items[:limit]
	}

	if hasMore && len(items) > 0 {
		lastItem := items[len(items)-1]
		nextCursor, _ = s.encodeCursor(&CursorData{
			ID:        lastItem.ID.String(),
			CreatedAt: lastItem.CreatedAt.Format(time.RFC3339Nano),
		})
	}

	return &models.PaginatedResponse{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Total:      total,
	}, nil
}

func (s *ItemService) SeedDatabase() error {
	var count int64
	s.db.Model(&models.Item{}).Count(&count)
	if count > 0 {
		return nil
	}

	sampleItems := []models.Item{
		{Name: "Laptop", Stock: 50, Price: 999.99},
		{Name: "Mouse", Stock: 200, Price: 25.99},
		{Name: "Keyboard", Stock: 150, Price: 75.50},
		{Name: "Monitor", Stock: 75, Price: 299.99},
		{Name: "Headphones", Stock: 100, Price: 149.99},
		{Name: "Webcam", Stock: 80, Price: 89.99},
		{Name: "USB Cable", Stock: 300, Price: 12.99},
		{Name: "Power Adapter", Stock: 120, Price: 45.00},
		{Name: "Tablet", Stock: 60, Price: 399.99},
		{Name: "Smartphone", Stock: 40, Price: 699.99},
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(sampleItems))

	for _, item := range sampleItems {
		wg.Add(1)
		go func(item models.Item) {
			defer wg.Done()
			if err := s.db.Create(&item).Error; err != nil {
				errors <- fmt.Errorf("failed to create item %s: %w", item.Name, err)
			}
		}(item)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ItemService) getFromCache(id string) *models.Item {
	if s.cache == nil {
		return nil
	}
	
	item, found := s.cache.Get(id)
	if !found {
		return nil
	}
	
	return item
}

func (s *ItemService) setCache(id string, item *models.Item) {
	if s.cache == nil {
		return
	}
	
	s.cache.SetWithTTL(id, item, 1, 5*time.Minute)
}

func (s *ItemService) invalidateCache() {
	if s.cache == nil {
		return
	}
	
	s.cache.Clear()
}

func (s *ItemService) Close() {
	if s.cache != nil {
		s.cache.Close()
	}
}

func (s *ItemService) encodeCursor(cursor *CursorData) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *ItemService) decodeCursor(cursor string) (*CursorData, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cursorData CursorData
	if err := json.Unmarshal(data, &cursorData); err != nil {
		return nil, err
	}

	return &cursorData, nil
}

func (s *ItemService) GetItemStats() (map[string]interface{}, error) {
	var stats struct {
		TotalItems    int64   `json:"total_items"`
		TotalValue    float64 `json:"total_value"`
		AveragePrice  float64 `json:"average_price"`
		LowStockItems int64   `json:"low_stock_items"`
	}

	if err := s.db.Model(&models.Item{}).Count(&stats.TotalItems).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Item{}).Select("SUM(price * stock) as total_value, AVG(price) as average_price").Scan(&stats).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Item{}).Where("stock < ?", 10).Count(&stats.LowStockItems).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_items":     stats.TotalItems,
		"total_value":     stats.TotalValue,
		"average_price":   stats.AveragePrice,
		"low_stock_items": stats.LowStockItems,
	}, nil
}
