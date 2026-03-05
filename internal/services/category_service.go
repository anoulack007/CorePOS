package services

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
)

type categoryService struct {
	repo ports.CategoryRepository
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(repo ports.CategoryRepository) ports.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories(storeID uint) ([]domain.Category, error) {
	return s.repo.FindAll(storeID)
}

func (s *categoryService) GetCategory(storeID, id uint) (*domain.Category, error) {
	return s.repo.FindByID(storeID, id)
}

func (s *categoryService) CreateCategory(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) UpdateCategory(category *domain.Category) error {
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(storeID, id uint) error {
	return s.repo.Delete(storeID, id)
}
