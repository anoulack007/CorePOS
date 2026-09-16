package services

import (
	"strings"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type inventoryService struct {
	db            *gorm.DB
	inventoryRepo ports.InventoryRepository
}

func NewInventoryService(db *gorm.DB, inventoryRepo ports.InventoryRepository) ports.InventoryService {
	return &inventoryService{db: db, inventoryRepo: inventoryRepo}
}

func (s *inventoryService) AdjustStock(storeID, productID uuid.UUID, userID *uuid.UUID, movementType string, quantityChanged int, notes string, evidenceURL string) error {
	movementType = strings.ToUpper(strings.TrimSpace(movementType))
	if !validMovementType(movementType) || quantityChanged == 0 {
		return domain.ErrInvalidInventoryMovement
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		var product domain.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("store_id = ? AND id = ?", storeID, productID).
			First(&product).Error; err != nil {
			return err
		}

		newQuantity := product.StockQuantity + quantityChanged
		if newQuantity < 0 {
			return domain.ErrInsufficientStock
		}
		if err := tx.Model(&domain.Product{}).
			Where("store_id = ? AND id = ?", storeID, productID).
			Update("stock_quantity", newQuantity).Error; err != nil {
			return err
		}

		movement := &domain.InventoryMovement{
			StoreID:         storeID,
			ProductID:       productID,
			UserID:          userID,
			MovementType:    movementType,
			QuantityChanged: quantityChanged,
			Notes:           notes,
			EvidenceURL:     evidenceURL,
		}

		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

func validMovementType(movementType string) bool {
	switch movementType {
	case "RESTOCK", "SALE", "RETURN", "DAMAGE", "ADJUSTMENT":
		return true
	default:
		return false
	}
}

func (s *inventoryService) GetStockHistory(storeID uuid.UUID, productID *uuid.UUID) ([]domain.InventoryMovement, error) {
	return s.inventoryRepo.GetHistory(storeID, productID)
}
