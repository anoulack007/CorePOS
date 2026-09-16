package services

import (
	"fmt"
	"strings"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
)

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) ports.ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(storeID uuid.UUID) ([]domain.Product, error) {
	return s.repo.FindAll(storeID)
}

func (s *productService) GetProduct(storeID, id uuid.UUID) (*domain.Product, error) {
	return s.repo.FindByID(storeID, id)
}

func (s *productService) CreateProduct(product *domain.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(storeID, id uuid.UUID) error {
	return s.repo.Delete(storeID, id)
}

func validateProduct(product *domain.Product) error {
	product.Name = strings.TrimSpace(product.Name)
	if product.Name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidProduct)
	}
	if product.Price < 0 || product.CostPrice < 0 {
		return fmt.Errorf("%w: prices cannot be negative", domain.ErrInvalidProduct)
	}
	if product.StockQuantity < 0 {
		return fmt.Errorf("%w: stock quantity cannot be negative", domain.ErrInvalidProduct)
	}
	return nil
}
