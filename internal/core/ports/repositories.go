package ports

import "github.com/anoulack007/core-pos/internal/core/domain"

// ProductRepository defines the data access interface for products.
type ProductRepository interface {
	FindAll(storeID uint) ([]domain.Product, error)
	FindByID(storeID, id uint) (*domain.Product, error)
	FindByBarcode(storeID uint, barcode string) (*domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(storeID, id uint) error
}

// CategoryRepository defines the data access interface for categories.
type CategoryRepository interface {
	FindAll(storeID uint) ([]domain.Category, error)
	FindByID(storeID, id uint) (*domain.Category, error)
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(storeID, id uint) error
}

// OrderRepository defines the data access interface for orders.
type OrderRepository interface {
	FindAll(storeID uint) ([]domain.Order, error)
	FindByID(storeID, id uint) (*domain.Order, error)
	Create(order *domain.Order) error
	UpdateStatus(storeID, id uint, status domain.OrderStatus) error
}

// UserRepository defines the data access interface for users.
type UserRepository interface {
	FindByID(storeID, id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

// PaymentRepository defines the data access interface for payments.
type PaymentRepository interface {
	Create(payment *domain.Payment) error
	FindByOrderID(orderID uint) ([]domain.Payment, error)
}
