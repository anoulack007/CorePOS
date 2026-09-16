package repositories

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(storeID, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	err := r.db.Where("store_id = ? AND id = ?", storeID, id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAllByStoreID(storeID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Where("store_id = ?", storeID).Order("created_at asc").Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) CreateInitialOwner(user *domain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var store domain.Store
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&store, "id = ?", user.StoreID).Error; err != nil {
			return err
		}

		var userCount int64
		if err := tx.Model(&domain.User{}).Where("store_id = ?", user.StoreID).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount > 0 {
			return domain.ErrStoreAlreadyInitialized
		}

		user.Role = domain.RoleOwner
		return tx.Create(user).Error
	})
}
