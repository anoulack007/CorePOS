package ports

import "github.com/anoulack007/core-pos/internal/core/domain"

// ProductService defines the business logic interface for products.
type ProductService interface {
	GetAllProducts(storeID uint) ([]domain.Product, error)
	GetProduct(storeID, id uint) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
	DeleteProduct(storeID, id uint) error
}

// OrderService defines the business logic interface for orders.
type OrderService interface {
	GetAllOrders(storeID uint) ([]domain.Order, error)
	GetOrder(storeID, id uint) (*domain.Order, error)
	CreateOrder(order *domain.Order) error
	VoidOrder(storeID, id uint) error
}

// CategoryService defines the business logic interface for categories.
type CategoryService interface {
	GetAllCategories(storeID uint) ([]domain.Category, error)
	GetCategory(storeID, id uint) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(storeID, id uint) error
}
