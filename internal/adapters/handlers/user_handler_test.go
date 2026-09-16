package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepositoryStub struct {
	created *domain.User
}

func (r *userRepositoryStub) FindByID(storeID, id uuid.UUID) (*domain.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *userRepositoryStub) FindByUsername(username string) (*domain.User, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *userRepositoryStub) FindAllByStoreID(storeID uuid.UUID) ([]domain.User, error) {
	return nil, nil
}
func (r *userRepositoryStub) Create(user *domain.User) error {
	r.created = user
	return nil
}
func (r *userRepositoryStub) CreateInitialOwner(user *domain.User) error { return nil }

func TestUserCreateAllowsOwnerToCreateAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storeID := uuid.New()
	repo := &userRepositoryStub{}
	handler := NewUserHandler(services.NewUserService(repo))
	router := gin.New()
	router.POST("/stores/:storeId/users", withRole(domain.RoleOwner), handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+storeID.String()+"/users", strings.NewReader(`{"username":"manager","password":"strong-password","role":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
	if repo.created == nil || repo.created.StoreID != storeID || repo.created.Role != domain.RoleAdmin {
		t.Fatalf("unexpected created staff: %#v", repo.created)
	}
	if strings.Contains(w.Body.String(), "password") {
		t.Fatalf("response exposed password data: %s", w.Body.String())
	}
}

func TestUserCreateBlocksAdminFromCreatingAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewUserHandler(services.NewUserService(&userRepositoryStub{}))
	router := gin.New()
	router.POST("/stores/:storeId/users", withRole(domain.RoleAdmin), handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+uuid.NewString()+"/users", strings.NewReader(`{"username":"manager","password":"strong-password","role":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}
}

func withRole(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.ContextRoleKey, role)
		c.Next()
	}
}
