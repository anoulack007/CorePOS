package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindAll(storeID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("store_id = ?", storeID).
		Preload("Items.Product").
		Preload("Payments").
		Order("created_at desc").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) FindByID(storeID, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Where("store_id = ? AND id = ?", storeID, id).
		Preload("Items.Product").
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

func (r *orderRepository) UpdateStatus(storeID, id uuid.UUID, status domain.OrderStatus) error {
	result := r.db.Model(&domain.Order{}).
		Where("store_id = ? AND id = ?", storeID, id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
