package main

import (
	"log"

	"github.com/anoulack007/core-pos/config"
	"github.com/anoulack007/core-pos/internal/adapters/handlers"
	"github.com/anoulack007/core-pos/internal/adapters/repositories"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ── Load Config ──────────────────────────────────────────────
	cfg := config.LoadConfig()

	// ── Connect Database ─────────────────────────────────────────
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")

	// ── Auto Migrate ─────────────────────────────────────────────
	err = db.AutoMigrate(
		&domain.Store{},
		&domain.SubscriptionHistory{},
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.InventoryMovement{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")

	// ── Repositories (Adapters) ──────────────────────────────────
	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)

	// ── Services (Business Logic) ────────────────────────────────
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	orderService := services.NewOrderService(orderRepo, productRepo, paymentRepo)

	// ── Handlers (HTTP Adapters) ─────────────────────────────────
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// ── Gin Router ───────────────────────────────────────────────
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes scoped by store
	api := r.Group("/api/v1/stores/:storeId")
	{
		// Products
		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", productHandler.Create)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}

		// Categories
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// Orders
		orders := api.Group("/orders")
		{
			orders.GET("", orderHandler.GetAll)
			orders.GET("/:id", orderHandler.GetByID)
			orders.POST("", orderHandler.Create)
			orders.PATCH("/:id/void", orderHandler.Void)
		}
	}

	// ── Start Server ─────────────────────────────────────────────
	log.Printf("CorePOS server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
