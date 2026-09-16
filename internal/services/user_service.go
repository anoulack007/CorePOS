package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetStaff(storeID uuid.UUID) ([]domain.User, error) {
	return s.repo.FindAllByStoreID(storeID)
}

func (s *userService) CreateStaff(actorRole domain.UserRole, user *domain.User, password string) error {
	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" || len(password) < 8 {
		return fmt.Errorf("%w: username and password of at least 8 characters are required", domain.ErrInvalidStaff)
	}

	if !canAssignRole(actorRole, user.Role) {
		return domain.ErrForbiddenRoleAssignment
	}

	if _, err := s.repo.FindByUsername(user.Username); err == nil {
		return domain.ErrUsernameAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrUsernameAlreadyExists
		}
		return err
	}
	return nil
}

func canAssignRole(actorRole, targetRole domain.UserRole) bool {
	switch actorRole {
	case domain.RoleOwner:
		return targetRole == domain.RoleAdmin || targetRole == domain.RoleCashier
	case domain.RoleAdmin:
		return targetRole == domain.RoleCashier
	default:
		return false
	}
}
