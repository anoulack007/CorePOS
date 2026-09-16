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
	CreateOrder(storeID, userID uuid.UUID, lines []OrderLineInput) (*domain.Order, error)
	VoidOrder(storeID, id, actorUserID uuid.UUID) error
}

type OrderLineInput struct {
	ProductID uuid.UUID
	Quantity  int
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

type UserService interface {
	GetStaff(storeID uuid.UUID) ([]domain.User, error)
	CreateStaff(actorRole domain.UserRole, user *domain.User, password string) error
}

type InventoryService interface {
	AdjustStock(storeID, productID uuid.UUID, userID *uuid.UUID, movementType string, quantityChanged int, notes string, evidenceURL string) error
	GetStockHistory(storeID uuid.UUID, productID *uuid.UUID) ([]domain.InventoryMovement, error)
}
