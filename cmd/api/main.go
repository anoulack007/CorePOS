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
		Logger:         gormlogger.Default.LogMode(gormlogger.Silent),
		TranslateError: true,
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
		&domain.InventoryMovement{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
		&domain.SubscriptionHistory{},
	)

	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	minioClient := config.ConnectMinIO(cfg)

	log.Println("✅ Database migrated successfully!")

	// Repositories
	productRepo := repositories.NewProductRepository(db)
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Services
	productService := services.NewProductService(productRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	categoryService := services.NewCategoryService(categoryRepo)
	inventoryService := services.NewInventoryService(db, inventoryRepo)
	userService := services.NewUserService(userRepo)
	orderService := services.NewOrderService(db, orderRepo)
	// Handlers
	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(db)
	authHandler := handlers.NewAuthHandler(authService, minioClient, cfg)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)

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

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", middleware.Auth(cfg.JWTSecret), authHandler.Logout)
	}

	// Store-scoped routes
	store := api.Group("/stores/:storeId")
	store.Use(middleware.Auth(cfg.JWTSecret), middleware.AuthorizeStoreAccess())
	{
		products := store.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), productHandler.Create)
			products.PUT("/:id", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), productHandler.Update)
			products.DELETE("/:id", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), productHandler.Delete)

		}

		categories := store.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), categoryHandler.Create)
			categories.PUT("/:id", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), categoryHandler.Update)
			categories.DELETE("/:id", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), categoryHandler.Delete)
		}

		inventory := store.Group("/inventory")
		{
			inventory.POST("/adjust", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), inventoryHandler.AdjustStock)
			inventory.GET("/history", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), inventoryHandler.GetHistory)
		}

		users := store.Group("/users")
		users.Use(middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin))
		{
			users.GET("", userHandler.GetAll)
			users.POST("", userHandler.Create)
		}

		orders := store.Group("/orders")
		{
			orders.GET("", orderHandler.GetAll)
			orders.GET("/:id", orderHandler.GetByID)
			orders.POST("", orderHandler.Create)
			orders.POST("/:id/void", middleware.RequireRoles(domain.RoleOwner, domain.RoleAdmin), orderHandler.Void)
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
