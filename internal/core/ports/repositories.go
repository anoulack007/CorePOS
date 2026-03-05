package ports

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/google/uuid"
)

type ProductRepository interface {
	FindAll(storeID uuid.UUID) ([]domain.Product, error)
	FindByID(storeID, id uuid.UUID) (*domain.Product,error)
	FindByBarcode(storeID uuid.UUID, barcode string) (*domain.Product,error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(storeId, id uuid.UUID) error
}

type CategoryRepository interface {
	FindAll(storeID uuid.UUID) ([]domain.Category,error)
	FindByID(storeID,id uuid.UUID) (*domain.Category, error)
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(storeID, id uuid.UUID) error
}

type OrderRepository interface {
	FindAll(storeID uuid.UUID) ([]domain.Order, error)
	FindByID(storeID, id uuid.UUID) (*domain.Order, error)
	Create(order *domain.Order) error
	UpdateStatus(storeID, id uuid.UUID, status domain.OrderStatus) error
}

type UserRepository interface {
	FindByID(storeID, id uuid.UUID) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}
type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindByOrderID(orderID uuid.UUID) ([]domain.Payment, error)
}