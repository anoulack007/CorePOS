package main

import (
	"fmt"
	"log"

	"github.com/anoulack007/core-pos/config"
	"github.com/anoulack007/core-pos/internal/adapters/handlers"
	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/adapters/repositories"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to DB
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully!")

	err = db.AutoMigrate(
		&domain.Store{},
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
	)

	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	log.Println("✅ Database migrated successfully!")

	// Repositories
	productRepo := repositories.NewProductRepository(db)
	// Services
	productService := services.NewProductService(productRepo)
	// Handlers
	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(db)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies(nil)

	// Middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.Security())
	r.Use(middleware.Compression())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// Store routes
	api.POST("/stores", storeHandler.Create)
	api.GET("/stores", storeHandler.GetAll)

	// Store-scoped routes
	store := api.Group("/stores/:storeId")
	{
		products := store.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", productHandler.Create)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}
	}

	// Start
	port := cfg.AppPort
	fmt.Printf("\n🚀 CorePOS API running on http://localhost:%s\n", port)
	fmt.Printf("📋 Health: http://localhost:%s/health\n\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
