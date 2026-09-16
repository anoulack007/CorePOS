package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ports.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(storeID uuid.UUID) ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Where("store_id = ?", storeID).Preload("Category").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(storeID, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("store_id = ? AND id = ?", storeID, id).Preload("Category").First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindByBarcode(storeID uuid.UUID, barcode string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("store_id = ? AND barcode = ?", storeID, barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *domain.Product) error {
	if err := r.validateCategory(product); err != nil {
		return err
	}
	return r.db.Omit("Store", "Category").Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	if err := r.validateCategory(product); err != nil {
		return err
	}
	return r.db.Model(&domain.Product{}).
		Where("store_id = ? AND id = ?", product.StoreID, product.ID).
		Select("category_id", "name", "barcode", "cost_price", "price", "stock_quantity", "image_url").
		Updates(product).Error
}

func (r *productRepository) Delete(storeID, id uuid.UUID) error {
	result := r.db.Where("store_id = ? AND id = ?", storeID, id).Delete(&domain.Product{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *productRepository) validateCategory(product *domain.Product) error {
	if product.CategoryID == nil {
		return nil
	}

	var count int64
	if err := r.db.Model(&domain.Category{}).
		Where("store_id = ? AND id = ?", product.StoreID, *product.CategoryID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrCategoryNotInStore
	}
	return nil
}
