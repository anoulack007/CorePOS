package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new GORM-backed OrderRepository.
func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll(storeID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("store_id = ?", storeID).
		Preload("Items").Preload("Items.Product").
		Preload("Payments").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(storeID, id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Where("store_id = ? AND id = ?", storeID, id).
		Preload("Items").Preload("Items.Product").
		Preload("Payments").
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) UpdateStatus(storeID, id uint, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).
		Where("store_id = ? AND id = ?", storeID, id).
		Update("status", status).Error
}
