# Inventory Management API

A high-performance inventory management system with CRUD operations, pagination, filtering, sorting, and rate limiting.

## ⚠️ Important Note

**Please wait 1-2 minutes after starting the server!** The database connection and migrations take time to complete. You'll see "Server starting on port 8080" when it's ready. The API will be available once the connection is established.

## 📝 Project Disclaimer

**This is a simple project to try out Go and Gin framework.** Everything is stored in PostgreSQL with basic caching and no advanced security features. This project was created for educational purposes to explore RESTful APIs with Go.

## Tech Stack
- **Go** ~> 1.23
- **Gin** Web Framework
- **PostgreSQL** Database
- **GORM** ORM
- **Ristretto** Cache
- **Docker** (optional)

## 🐳 Docker (Recommended)

### Basic Version (Local Access)
```bash
# Build the Docker image
docker build -t inventory-api .

# Run with Docker Compose (includes PostgreSQL)
docker-compose up -d
```

Access at: **http://localhost:8080**

**Services included:**
- **API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **pgAdmin**: http://localhost:5050 (admin@inventory.com / admin)

### Docker Commands
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Clean up (removes volumes)
docker-compose down -v
```

## 💻 Manual Installation

### Requirements
- **Go** ~> 1.23
- **PostgreSQL** 12 or higher
- **Git** (for cloning)

### Setup
```bash
git clone <repository-url>
cd inventory-api

# Install dependencies
go mod tidy

# Create database
createdb inventory_db

# Run the application
go run main.go
```

Access at: **http://localhost:8080**

### Version Check Commands
```bash
# Check Go version
go version

# Check PostgreSQL version
psql --version
```

## 🚀 Features

- **CRUD Operations**: Create, Read, Update, Delete inventory items
- **Pagination**: Cursor-based pagination for large datasets
- **Filtering**: Filter by name, stock levels, and price ranges
- **Sorting**: Sort by name, stock, price, or creation date
- **Rate Limiting**: 1 request per second with burst capacity of 5
- **Caching**: High-performance Ristretto cache
- **API Documentation**: Swagger UI at `/api/v1/swagger/index.html`

## 📸 API Endpoints

- `GET /api/v1/inventory` - List all items (with pagination, filtering, sorting)
- `POST /api/v1/inventory` - Create new item
- `GET /api/v1/inventory/:id` - Get item by ID
- `PUT /api/v1/inventory/:id` - Update item
- `DELETE /api/v1/inventory/:id` - Delete item
- `GET /api/v1/inventory/stats` - Get inventory statistics
- `POST /api/v1/inventory/seed` - Seed database with sample data
- `GET /health` - Health check
- `GET /api/v1/swagger/index.html` - API documentation

## 🔧 Troubleshooting

### Database Connection Issues
- Wait 1-2 minutes after starting the server
- Check logs for "Server starting on port 8080"
- Ensure PostgreSQL is running: `pg_isready`
- Check database exists: `psql -l | grep inventory_db`

### Docker Issues
- Ensure ports 8080 and 5432 are not in use: `docker ps`
- Stop conflicting containers: `docker stop $(docker ps -q)`

### Manual Installation Issues
- Ensure Go 1.23+ is installed: `go version`
- Check PostgreSQL is running: `sudo systemctl status postgresql`
- Verify database exists: `psql -l | grep inventory_db`

## 📦 File Sizes

- **Docker Image**: ~50MB (compressed)
- **Source Code**: ~2MB
- **Dependencies**: ~100MB (first build only)

## 🎯 Quick Start

**Fastest way to get started:**
```bash
docker-compose up -d
```

Then open **http://localhost:8080** and wait 1-2 minutes for the database connection to establish!

**API Documentation**: http://localhost:8080/api/v1/swagger/index.html