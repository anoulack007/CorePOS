package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) ports.InventoryRepository {
	return &inventoryRepository{db:db}
}

func (r *inventoryRepository) LogMovement(movement *domain.InventoryMovement) error {
	return r.db.Create(movement).Error
}


func (r *inventoryRepository) GetHistory(storeID uuid.UUID, productID *uuid.UUID) ([]domain.InventoryMovement,error){
	var movements []domain.InventoryMovement
	query := r.db.Where("store_id = ?", storeID)

	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}

	err := query.Order("created_at desc").Preload("User").Preload("Product").Find(&movements).Error

	if err != nil {
		return nil,err
	}

	return movements,err
}