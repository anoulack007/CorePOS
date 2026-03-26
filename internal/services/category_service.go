package services

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
)

type categoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) ports.CategoryService{
	return &categoryService{repo:repo}
}

func(s *categoryService) GetAllCategories(storeID uuid.UUID) ([]domain.Category, error) {
	return s.repo.FindAll(storeID)
}

func(s *categoryService) GetCategory(storeID, id uuid.UUID) (*domain.Category, error) {
	return s.repo.FindByID(storeID,id)
}


func(s *categoryService) CreateCategory(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *categoryService) UpdateCategory(category *domain.Category) error {
	return s.repo.Update(category)
}

func(s *categoryService) DeleteCategory(storeID, id uuid.UUID) error {
	return s.repo.Delete(storeID, id)
}