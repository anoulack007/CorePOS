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

type productRepositoryStub struct {
	created *domain.Product
}

func (r *productRepositoryStub) FindAll(storeID uuid.UUID) ([]domain.Product, error) {
	return nil, nil
}
func (r *productRepositoryStub) FindByID(storeID, id uuid.UUID) (*domain.Product, error) {
	return nil, nil
}
func (r *productRepositoryStub) FindByBarcode(storeID uuid.UUID, barcode string) (*domain.Product, error) {
	return nil, nil
}
func (r *productRepositoryStub) Create(product *domain.Product) error {
	r.created = product
	return nil
}
func (r *productRepositoryStub) Update(product *domain.Product) error { return nil }
func (r *productRepositoryStub) Delete(storeID, id uuid.UUID) error   { return nil }

func TestProductCreateUsesProductionHandlerAndService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storeID := uuid.New()
	repo := &productRepositoryStub{}
	handler := NewProductHandler(services.NewProductService(repo))
	router := gin.New()
	router.POST("/stores/:storeId/products", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+storeID.String()+"/products", strings.NewReader(`{"name":"Coffee","price":45,"stock_quantity":2}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
	if repo.created == nil || repo.created.StoreID != storeID || repo.created.Name != "Coffee" {
		t.Fatalf("product was not passed through the production flow: %#v", repo.created)
	}
}

func TestProductCreateRejectsInvalidProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(services.NewProductService(&productRepositoryStub{}))
	router := gin.New()
	router.POST("/stores/:storeId/products", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+uuid.NewString()+"/products", strings.NewReader(`{"name":"","price":-1}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}
