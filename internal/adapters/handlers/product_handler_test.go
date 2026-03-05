package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --- Mock ProductService ---

type mockProductService struct {
	products []domain.Product
	err      error
}

func (m *mockProductService) GetAllProducts(storeID uuid.UUID) ([]domain.Product, error) {
	return m.products, m.err
}

func (m *mockProductService) GetProduct(storeID, id uuid.UUID) (*domain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.Product{
		ID:      id,
		StoreID: storeID,
		Name:    gofakeit.ProductName(),
		Price:   gofakeit.Price(10, 500),
	}, nil
}

func (m *mockProductService) CreateProduct(product *domain.Product) error {
	return m.err
}

func (m *mockProductService) UpdateProduct(product *domain.Product) error {
	return m.err
}

func (m *mockProductService) DeleteProduct(storeID, id uuid.UUID) error {
	return m.err
}

// --- Helper ---

func setupRouter(handler *ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := r.Group("/stores/:storeId")
	{
		store.GET("/products", handler.GetAll)
		store.GET("/products/:id", handler.GetByID)
		store.POST("/products", handler.Create)
		store.PUT("/products/:id", handler.Update)
		store.DELETE("/products/:id", handler.Delete)
	}
	return r
}

func fakeStoreID() string   { return uuid.New().String() }
func fakeProductID() string { return uuid.New().String() }

// --- GetAll ---

func TestGetAll_200(t *testing.T) {
	fakeProducts := []domain.Product{
		{ID: uuid.New(), Name: gofakeit.ProductName(), Price: gofakeit.Price(10, 200)},
		{ID: uuid.New(), Name: gofakeit.ProductName(), Price: gofakeit.Price(10, 200)},
	}
	h := NewProductHandler(&mockProductService{products: fakeProducts})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stores/"+fakeStoreID()+"/products", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetAll_400(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stores/invalid-uuid/products", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GetByID ---

func TestGetByID_200(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stores/"+fakeStoreID()+"/products/"+fakeProductID(), nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetByID_400(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stores/"+fakeStoreID()+"/products/bad-id", nil)
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Create ---

func TestCreate_201(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	body := fmt.Sprintf(`{"name":"%s","barcode":"%s","price":%.2f,"stock_quantity":%d}`,
		gofakeit.ProductName(),
		gofakeit.DigitN(13),
		gofakeit.Price(10, 500),
		gofakeit.Number(1, 1000),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/stores/"+fakeStoreID()+"/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreate_400(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/stores/"+fakeStoreID()+"/products", strings.NewReader(`{invalid}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Update ---

func TestUpdate_200(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	body := fmt.Sprintf(`{"name":"%s","price":%.2f}`, gofakeit.ProductName(), gofakeit.Price(10, 500))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/stores/"+fakeStoreID()+"/products/"+fakeProductID(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_400(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/stores/"+fakeStoreID()+"/products/"+fakeProductID(), strings.NewReader(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Delete ---

func TestDelete_200(t *testing.T) {
	h := NewProductHandler(&mockProductService{})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/stores/"+fakeStoreID()+"/products/"+fakeProductID(), nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDelete_400(t *testing.T) {
	h := NewProductHandler(&mockProductService{err: errors.New("delete failed")})
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/stores/"+fakeStoreID()+"/products/"+fakeProductID(), nil)
	r.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
