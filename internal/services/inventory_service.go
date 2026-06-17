package services

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryService struct {
	db            *gorm.DB
	inventoryRepo ports.InventoryRepository
	productRepo   ports.ProductRepository
}

func NewInventoryService(db *gorm.DB, inventoryRepo ports.InventoryRepository, productRepo ports.ProductRepository) ports.InventoryService {
	return &inventoryService{db: db, inventoryRepo: inventoryRepo, productRepo: productRepo}
}

func (s *inventoryService) AdjustStock(storeID, productID uuid.UUID, userID *uuid.UUID, movementType string, quantityChanged int, notes string, evidenceURL string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch the product
		product, err := s.productRepo.FindByID(storeID, productID)
		if err != nil {
			return err
		}

		// 2. Adjust stock quantity
		product.StockQuantity += quantityChanged

		// 3. Save updated product
		if err := s.productRepo.Update(product); err != nil {
			return err
		}

		// 4. Log the movement
		movement := &domain.InventoryMovement{
			StoreID:         storeID,
			ProductID:       productID,
			UserID:          userID,
			MovementType:    movementType,
			QuantityChanged: quantityChanged,
			Notes:           notes,
			EvidenceURL:     evidenceURL,
		}

		if err := s.inventoryRepo.LogMovement(movement); err != nil {
			return err
		}

		return nil
	})
}

func (s *inventoryService) GetStockHistory(storeID uuid.UUID, productID *uuid.UUID) ([]domain.InventoryMovement, error) {
	return s.inventoryRepo.GetHistory(storeID, productID)
}
