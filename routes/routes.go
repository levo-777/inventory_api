package routes

import (
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	"inventory-api/controllers"
	"inventory-api/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(cfg *utils.Config) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(utils.CORSMiddleware())

	// Apply rate limiting only to API routes, not to Swagger or health endpoints
	apiGroup := router.Group("/api")
	apiGroup.Use(utils.RateLimitMiddleware(cfg.RateLimit.Requests, cfg.RateLimit.Burst))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check database health
		if err := utils.Health(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  err.Error(),
			})
			return
		}

		// Get system info
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
			"system": gin.H{
				"goroutines": runtime.NumGoroutine(),
				"memory_mb":  m.Alloc / 1024 / 1024,
				"gc_runs":    m.NumGC,
			},
		})
	})

	// Swagger documentation (no rate limiting)
	router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes (with rate limiting)
	v1 := apiGroup.Group("/v1")
	{
		inventory := v1.Group("/inventory")
		{
			itemController := controllers.NewItemController()

			inventory.GET("", itemController.GetItems)
			inventory.POST("", itemController.CreateItem)
			inventory.GET("/stats", itemController.GetItemStats)
			inventory.POST("/seed", itemController.SeedDatabase)
			inventory.GET("/:id", itemController.GetItem)
			inventory.PUT("/:id", itemController.UpdateItem)
			inventory.DELETE("/:id", itemController.DeleteItem)
		}
	}

	// Profiling endpoints (available in all modes for development)
	debug := router.Group("/debug")
	{
		debug.GET("/pprof/", gin.WrapF(http.HandlerFunc(pprof.Index)))
		debug.GET("/pprof/cmdline", gin.WrapF(http.HandlerFunc(pprof.Cmdline)))
		debug.GET("/pprof/profile", gin.WrapF(http.HandlerFunc(pprof.Profile)))
		debug.GET("/pprof/symbol", gin.WrapF(http.HandlerFunc(pprof.Symbol)))
		debug.GET("/pprof/trace", gin.WrapF(http.HandlerFunc(pprof.Trace)))
		debug.GET("/pprof/goroutine", gin.WrapF(http.HandlerFunc(pprof.Handler("goroutine").ServeHTTP)))
		debug.GET("/pprof/heap", gin.WrapF(http.HandlerFunc(pprof.Handler("heap").ServeHTTP)))
		debug.GET("/pprof/block", gin.WrapF(http.HandlerFunc(pprof.Handler("block").ServeHTTP)))
		debug.GET("/pprof/mutex", gin.WrapF(http.HandlerFunc(pprof.Handler("mutex").ServeHTTP)))
		debug.GET("/pprof/allocs", gin.WrapF(http.HandlerFunc(pprof.Handler("allocs").ServeHTTP)))
	}

	return router
}
