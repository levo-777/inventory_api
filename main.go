package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-api/routes"
	"inventory-api/utils"

	_ "inventory-api/docs"

	"github.com/gin-gonic/gin"
)

// @title Inventory Management API
// @version 1.0
// @description A comprehensive inventory management system with CRUD operations, pagination, filtering, sorting, and rate limiting
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	cfg, err := utils.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := utils.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer utils.Close()

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env == "development" {
		utils.Info.Println("Running AutoMigrate in development mode...")
		if err := utils.Migrate(); err != nil {
			utils.Error.Printf("Failed to migrate database: %v", err)
		}
	}

	itemService := utils.NewItemService()
	if err := itemService.SeedDatabase(); err != nil {
		utils.Error.Printf("Failed to seed database: %v", err)
	}

	router := routes.SetupRoutes(cfg)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		utils.Info.Printf("Server starting on port %s", cfg.Server.Port)
		utils.Info.Printf("API Documentation available at: http://localhost:%s/api/v1/swagger/index.html", cfg.Server.Port)
		utils.Info.Printf("Health check available at: http://localhost:%s/health", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Error.Printf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.Error.Printf("Server forced to shutdown: %v", err)
	}

	utils.Info.Println("Server exited")
}
