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

	minioClient := config.ConnectMinIO(cfg)

	log.Println("✅ Database migrated successfully!")

	// Repositories
	productRepo := repositories.NewProductRepository(db)
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	// Services
	productService := services.NewProductService(productRepo)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	categoryService := services.NewCategoryService(categoryRepo)
	// Handlers
	productHandler := handlers.NewProductHandler(productService)
	storeHandler := handlers.NewStoreHandler(db)
	authHandler := handlers.NewAuthHandler(authService, minioClient, cfg)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

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
			products.POST("", middleware.RequireRoles("owner", "admin"), productHandler.Create)
			products.PUT("/:id", middleware.RequireRoles("owner", "admin"), productHandler.Update)
			products.DELETE("/:id", middleware.RequireRoles("owner", "admin"), productHandler.Delete)
		}

		categories := store.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", middleware.RequireRoles("owner", "admin"), categoryHandler.Create)
			categories.PUT("/:id", middleware.RequireRoles("owner", "admin"), categoryHandler.Update)
			categories.DELETE("/:id", middleware.RequireRoles("owner", "admin"), categoryHandler.Delete)
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
