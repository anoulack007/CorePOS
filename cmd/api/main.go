package main

import (
	"log"

	"github.com/anoulack007/core-pos/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to DB
	db, err := gorm.Open(postgres.Open(cfg.DSN()),&gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

}
