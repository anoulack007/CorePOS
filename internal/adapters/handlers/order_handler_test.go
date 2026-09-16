package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type orderServiceStub struct {
	storeID uuid.UUID
	userID  uuid.UUID
	lines   []ports.OrderLineInput
}

func (s *orderServiceStub) GetAllOrders(storeID uuid.UUID) ([]domain.Order, error) {
	return nil, nil
}
func (s *orderServiceStub) GetOrder(storeID, id uuid.UUID) (*domain.Order, error) {
	return nil, nil
}
func (s *orderServiceStub) CreateOrder(storeID, userID uuid.UUID, lines []ports.OrderLineInput) (*domain.Order, error) {
	s.storeID = storeID
	s.userID = userID
	s.lines = lines
	return &domain.Order{ID: uuid.New(), StoreID: storeID, UserID: userID, TotalAmount: 20}, nil
}
func (s *orderServiceStub) VoidOrder(storeID, id, actorUserID uuid.UUID) error { return nil }

func TestOrderCreateUsesAuthenticatedUserAndRequestLines(t *testing.T) {
	gin.SetMode(gin.TestMode)
	storeID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()
	service := &orderServiceStub{}
	handler := NewOrderHandler(service)
	router := gin.New()
	router.POST("/stores/:storeId/orders", func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, userID)
		c.Next()
	}, handler.Create)

	body := `{"total_amount":1,"items":[{"product_id":"` + productID.String() + `","quantity":2,"unit_price":1}]}`
	req := httptest.NewRequest(http.MethodPost, "/stores/"+storeID.String()+"/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}
	if service.storeID != storeID || service.userID != userID {
		t.Fatalf("unexpected order identity: store=%s user=%s", service.storeID, service.userID)
	}
	if len(service.lines) != 1 || service.lines[0].ProductID != productID || service.lines[0].Quantity != 2 {
		t.Fatalf("unexpected order lines: %#v", service.lines)
	}
}

func TestOrderCreateRejectsEmptyItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewOrderHandler(&orderServiceStub{})
	router := gin.New()
	router.POST("/stores/:storeId/orders", func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, uuid.New())
		c.Next()
	}, handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores/"+uuid.NewString()+"/orders", strings.NewReader(`{"items":[]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}
