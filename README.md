# Inventory Management API

A **minimalistic** and lightweight inventory management system API built with Go, Gin, and PostgreSQL. Perfect for small to medium businesses that need a simple, fast, and reliable inventory tracking solution.

## 🚀 Quick Overview

This API provides a clean and efficient way to manage inventory through RESTful endpoints. It supports comprehensive inventory operations with advanced features:

- **Items**: Product inventory with stock tracking
- **CRUD Operations**: Full create, read, update, delete functionality
- **Advanced Features**: Pagination, filtering, sorting, and rate limiting
- **Performance**: High-performance caching with Ristretto
- **Monitoring**: Built-in health checks and profiling

## 🛠 Tech Stack

- **Backend**: Go 1.23
- **Web Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Cache**: Ristretto (high-performance)
- **Testing**: Go testing framework with SQLite
- **Containerization**: Docker & Docker Compose
- **Documentation**: Swagger/OpenAPI

## API Endpoints

### Items
- `GET /api/v1/inventory` - Get all items (with pagination, filtering, sorting)
- `GET /api/v1/inventory/:id` - Get item by ID
- `POST /api/v1/inventory` - Create new item
- `PUT /api/v1/inventory/:id` - Update item
- `DELETE /api/v1/inventory/:id` - Delete item
- `GET /api/v1/inventory/stats` - Get inventory statistics
- `POST /api/v1/inventory/seed` - Seed database with sample data

### System
- `GET /health` - Health check endpoint
- `GET /api/v1/swagger/index.html` - API documentation
- `GET /debug/pprof/*` - Performance profiling

## Data Models

### Item
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Smartphone",
  "stock": 40,
  "price": 699.99,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create Item Request
```json
{
  "name": "Laptop",
  "stock": 50,
  "price": 999.99
}
```

### Update Item Request
```json
{
  "name": "Gaming Laptop",
  "stock": 25,
  "price": 1299.99
}
```

### Paginated Response
```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "eyJpZCI6IjU1MGU4NDAwLWUyOWItNDFkNC1hNzE2LTQ0NjY1NTQ0MDAwMCIsImNyZWF0ZWRfYXQiOiIyMDI0LTAxLTAxVDAwOjAwOjAwWiJ9",
    "has_more": true,
    "total": 100
  }
}
```

## 📋 Prerequisites

### For Local Development (without Docker)
- **Go 1.23** or higher
- **PostgreSQL 15** or higher
- Git

### For Docker Deployment
- **Docker** 20.10+
- **Docker Compose** 2.0+

## 🚀 Getting Started

Choose your preferred setup method:

- [**Local Development**](#-local-development-setup) - Run directly on your machine
- [**Docker Development**](#-docker-development-setup) - Run with Docker Compose

---

## 💻 Local Development Setup

### 1. Clone the Repository
```bash
git clone https://github.com/levo-777/inventory_api.git
cd inventory_api
```

### 2. Install Go Dependencies
```bash
go mod download
```

### 3. Environment Configuration
```bash
# Copy environment template
cp env.example .env

# Edit with your database settings
nano .env
```

**Required environment variables:**
```env
ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inventory_db
SERVER_PORT=8080
RATE_LIMIT_REQUESTS=1
RATE_LIMIT_BURST=5
```

### 4. Database Setup
```bash
# Start PostgreSQL service
sudo systemctl start postgresql

# Create database
psql -U postgres -c "CREATE DATABASE inventory_db;"
```

### 5. Run the Application
```bash
# Start the API server
go run main.go
```

✅ **API available at:** `http://localhost:8080`

---

## 🐳 Docker Development Setup

### Quick Start (Recommended)

```bash
# Clone the repository
git clone https://github.com/levo-777/inventory_api.git
cd inventory_api

# Start all services with Docker Compose
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f
```

✅ **API available at:** `http://localhost:8080`
✅ **pgAdmin available at:** `http://localhost:5050`

## 📖 API Usage Examples

### Create an Item
```bash
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gaming Laptop",
    "stock": 25,
    "price": 1299.99
  }'
```

### Get All Items (with pagination)
```bash
curl "http://localhost:8080/api/v1/inventory?limit=10&sort=name&order=asc"
```

### Get Items with Filtering
```bash
# Filter by name
curl "http://localhost:8080/api/v1/inventory?name=laptop"

# Filter by minimum stock
curl "http://localhost:8080/api/v1/inventory?min_stock=50"
```

### Update an Item
```bash
curl -X PUT http://localhost:8080/api/v1/inventory/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Laptop",
    "stock": 30,
    "price": 1199.99
  }'
```

### Get Item by ID
```bash
curl http://localhost:8080/api/v1/inventory/550e8400-e29b-41d4-a716-446655440000
```

### Get Inventory Statistics
```bash
curl http://localhost:8080/api/v1/inventory/stats
```

### Seed Database with Sample Data
```bash
curl -X POST http://localhost:8080/api/v1/inventory/seed
```

### Health Check
```bash
curl http://localhost:8080/health
```

## 🔧 Advanced Features

### Pagination
- **Cursor-based pagination** for efficient large dataset handling
- Use `limit` parameter to control page size
- Use `cursor` parameter for next page navigation

### Filtering
- **By name**: `?name=keyword`
- **By minimum stock**: `?min_stock=50`

### Sorting
- **Sort by**: `name`, `stock`, `price`, `created_at`
- **Order**: `asc` or `desc`
- Example: `?sort=price&order=desc`

### Rate Limiting
- **1 request per second** with burst capacity of 5
- Applied to all API endpoints except health and documentation

### Caching
- **High-performance Ristretto cache** for frequently accessed items
- Automatic cache invalidation on updates
- 5-minute TTL for cached items

## 🧪 Testing

### Run Unit Tests
```bash
go test ./...
```

### Run Integration Tests
```bash
go test ./test/integrations/...
```

### Run with Coverage
```bash
go test -cover ./...
```

## 📊 Monitoring & Profiling

### Health Check
```bash
curl http://localhost:8080/health
```

### Performance Profiling
```bash
# CPU profile
curl http://localhost:8080/debug/pprof/profile

# Memory profile
curl http://localhost:8080/debug/pprof/heap

# Goroutine profile
curl http://localhost:8080/debug/pprof/goroutine
```

### API Documentation
Visit `http://localhost:8080/api/v1/swagger/index.html` for interactive API documentation.

## 🚀 Production Deployment

### Docker Production Build
```bash
# Build production image
docker build -t inventory-api .

# Run with production environment
docker run -d \
  --name inventory-api \
  -p 8080:8080 \
  --env-file docker.env \
  inventory-api
```

### Environment Variables for Production
```env
ENV=production
GIN_MODE=release
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-secure-password
DB_NAME=inventory_db
SERVER_PORT=8080
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_BURST=20
```

## 📁 Project Structure

```
inventory-api/
├── controllers/          # HTTP handlers
├── models/              # Data models and DTOs
├── routes/              # Route definitions
├── utils/               # Utilities (config, database, cache, etc.)
├── migrations/          # Database migrations
├── test/                # Test files
│   └── integrations/    # Integration tests
├── docs/                # Generated Swagger documentation
├── main.go              # Application entry point
├── docker-compose.yml   # Docker services
├── Dockerfile           # Container definition
└── README.md           # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Failed**
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Verify database exists
psql -U postgres -l | grep inventory_db
```

**Port Already in Use**
```bash
# Kill process using port 8080
sudo lsof -ti:8080 | xargs kill -9

# Or use a different port
SERVER_PORT=8081 go run main.go
```

**Rate Limit Exceeded**
- Wait 1 second between requests
- Use burst capacity (up to 5 requests quickly)
- Check rate limit headers in response

### Getting Help

- Check the [API Documentation](http://localhost:8080/api/v1/swagger/index.html)
- Review the health check endpoint: `GET /health`
- Check application logs for detailed error messages