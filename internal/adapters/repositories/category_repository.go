package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new GORM-backed CategoryRepository.
func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(storeID uint) ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Where("store_id = ?", storeID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(storeID, id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("store_id = ? AND id = ?", storeID, id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(storeID, id uint) error {
	return r.db.Where("store_id = ? AND id = ?", storeID, id).Delete(&domain.Category{}).Error
}
