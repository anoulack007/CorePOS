package services

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
)

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) ports.ProductService{
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(storeID uuid.UUID) ([]domain.Product, error){
	return s.repo.FindAll(storeID)
}

func (s *productService) GetProduct(storeID, id uuid.UUID) (*domain.Product, error){
	return s.repo.FindByID(storeID,id)
}

func (s *productService) CreateProduct(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(storeID, id uuid.UUID) error {
	return s.repo.Delete(storeID, id)
}