package ports

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAllProducts(storeID uuid.UUID) ([]domain.Product, error)
	GetProduct(storeID, id uuid.UUID) (*domain.Product, error)
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
	DeleteProduct(storeID, id uuid.UUID) error
}

type OrderService interface {
	GetAllOrders(storeID uuid.UUID) ([]domain.Order, error)
	GetOrder(storeID, id uuid.UUID) (*domain.Order, error)
	CreateOrder(order *domain.Order) error
	VoidOrder(storeID, id uuid.UUID) error
}

type CategoryService interface {
	GetAllCategories(storeID uuid.UUID) ([]domain.Category, error)
	GetCategory(storeID, id uuid.UUID) (*domain.Category, error)
	CreateCategory(category *domain.Category) error
	UpdateCategory(category *domain.Category) error
	DeleteCategory(storeID, id uuid.UUID) error
}

type AuthService interface {
	Register(user *domain.User, password string) error
	Login(username, password string) (string, string, error) // returns JWT token
	RefreshToken(token string) (newAccessToken string, newRefreshToken string, err error)
	Logout() error
}
