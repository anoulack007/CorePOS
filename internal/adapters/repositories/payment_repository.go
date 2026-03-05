package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new GORM-backed PaymentRepository.
func NewPaymentRepository(db *gorm.DB) ports.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByOrderID(orderID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Where("order_id = ?", orderID).Find(&payments).Error
	return payments, err
}
