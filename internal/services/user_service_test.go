package services

import (
	"errors"
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type staffRepositoryStub struct {
	created *domain.User
	users   []domain.User
}

func (r *staffRepositoryStub) FindByID(storeID, id uuid.UUID) (*domain.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *staffRepositoryStub) FindByUsername(username string) (*domain.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *staffRepositoryStub) FindAllByStoreID(storeID uuid.UUID) ([]domain.User, error) {
	return r.users, nil
}
func (r *staffRepositoryStub) Create(user *domain.User) error {
	r.created = user
	return nil
}
func (r *staffRepositoryStub) CreateInitialOwner(user *domain.User) error { return nil }

func TestOwnerCanCreateAdmin(t *testing.T) {
	repo := &staffRepositoryStub{}
	service := NewUserService(repo)
	user := &domain.User{StoreID: uuid.New(), Username: "manager", Role: domain.RoleAdmin}

	if err := service.CreateStaff(domain.RoleOwner, user, "strong-password"); err != nil {
		t.Fatalf("CreateStaff() error = %v", err)
	}
	if repo.created == nil || repo.created.Role != domain.RoleAdmin {
		t.Fatalf("expected admin to be created, got %#v", repo.created)
	}
	if bcrypt.CompareHashAndPassword([]byte(repo.created.PasswordHash), []byte("strong-password")) != nil {
		t.Fatal("expected staff password to be stored as a bcrypt hash")
	}
}

func TestAdminCannotCreateAdmin(t *testing.T) {
	service := NewUserService(&staffRepositoryStub{})
	user := &domain.User{Username: "another-admin", Role: domain.RoleAdmin}

	err := service.CreateStaff(domain.RoleAdmin, user, "strong-password")
	if !errors.Is(err, domain.ErrForbiddenRoleAssignment) {
		t.Fatalf("expected ErrForbiddenRoleAssignment, got %v", err)
	}
}

func TestAdminCanCreateCashier(t *testing.T) {
	repo := &staffRepositoryStub{}
	service := NewUserService(repo)
	user := &domain.User{Username: "cashier", Role: domain.RoleCashier}

	if err := service.CreateStaff(domain.RoleAdmin, user, "strong-password"); err != nil {
		t.Fatalf("CreateStaff() error = %v", err)
	}
	if repo.created == nil || repo.created.Role != domain.RoleCashier {
		t.Fatalf("expected cashier to be created, got %#v", repo.created)
	}
}

func TestStaffCannotCreateOwner(t *testing.T) {
	service := NewUserService(&staffRepositoryStub{})
	user := &domain.User{Username: "owner-two", Role: domain.RoleOwner}

	err := service.CreateStaff(domain.RoleOwner, user, "strong-password")
	if !errors.Is(err, domain.ErrForbiddenRoleAssignment) {
		t.Fatalf("expected owner assignment to be rejected, got %v", err)
	}
}
