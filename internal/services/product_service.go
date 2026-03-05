package services

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
)

type productService struct {
	repo ports.ProductRepository
}

// NewProductService creates a new ProductService.
func NewProductService(repo ports.ProductRepository) ports.ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(storeID uint) ([]domain.Product, error) {
	return s.repo.FindAll(storeID)
}

func (s *productService) GetProduct(storeID, id uint) (*domain.Product, error) {
	return s.repo.FindByID(storeID, id)
}

func (s *productService) CreateProduct(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(storeID, id uint) error {
	return s.repo.Delete(storeID, id)
}
