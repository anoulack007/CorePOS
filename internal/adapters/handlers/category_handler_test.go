package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type categoryRepositoryStub struct {
	created *domain.Category
}

func (r *categoryRepositoryStub) FindAll(storeID uuid.UUID) ([]domain.Category, error) {
	return nil, nil
}
func (r *categoryRepositoryStub) FindByID(storeID, id uuid.UUID) (*domain.Category, error) {
	return nil, nil
}
func (r *categoryRepositoryStub) Create(category *domain.Category) error {
	r.created = category
	return nil
}
func (r *categoryRepositoryStub) Update(category *domain.Category) error { return nil }
func (r *categoryRepositoryStub) Delete(storeID, id uuid.UUID) error     { return nil }

func TestCategoryCreateUsesProductionHandlerAndService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storeID := uuid.New()
	repo := &categoryRepositoryStub{}
	handler := NewCategoryHandler(services.NewCategoryService(repo))
	router := gin.New()
	router.POST("/stores/:storeId/categories", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+storeID.String()+"/categories", strings.NewReader(`{"name":"Drinks"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
	if repo.created == nil || repo.created.StoreID != storeID || repo.created.Name != "Drinks" {
		t.Fatalf("category was not passed through the production flow: %#v", repo.created)
	}
}

func TestCategoryCreateRejectsBlankName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewCategoryHandler(services.NewCategoryService(&categoryRepositoryStub{}))
	router := gin.New()
	router.POST("/stores/:storeId/categories", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+uuid.NewString()+"/categories", strings.NewReader(`{"name":"  "}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}
