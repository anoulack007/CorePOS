package services

import (
	"github.com/anoulack007/core-pos/internal/core/ports"
	"gorm.io/gorm"
)

type inventoryService struct {
	db            *gorm.DB
	inventoryRepo ports.InventoryRepository
	productRepo   ports.ProductRepository
}

func NewInventoryService(db *gorm.DB, inventoryRepo ports.InventoryRepository, productRepo ports.ProductRepository) *inventoryService {
	return &inventoryService{db: db, inventoryRepo: inventoryRepo, productRepo: productRepo}
}
