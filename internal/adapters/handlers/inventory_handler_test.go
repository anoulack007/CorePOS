package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type inventoryServiceStub struct {
	userID *uuid.UUID
}

func (s *inventoryServiceStub) AdjustStock(storeID, productID uuid.UUID, userID *uuid.UUID, movementType string, quantityChanged int, notes string, evidenceURL string) error {
	s.userID = userID
	return nil
}

func (s *inventoryServiceStub) GetStockHistory(storeID uuid.UUID, productID *uuid.UUID) ([]domain.InventoryMovement, error) {
	return nil, nil
}

func TestAdjustStockUsesAuthenticatedUserForAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authenticatedUserID := uuid.New()
	service := &inventoryServiceStub{}
	handler := NewInventoryHandler(service)
	router := gin.New()
	router.POST("/stores/:storeId/inventory/adjust", func(c *gin.Context) {
		c.Set(middleware.ContextUserIDKey, authenticatedUserID)
		c.Next()
	}, handler.AdjustStock)

	body := `{"product_id":"` + uuid.NewString() + `","user_id":"` + uuid.NewString() + `","movement_type":"RESTOCK","quantity_changed":5}`
	req := httptest.NewRequest(http.MethodPost, "/stores/"+uuid.NewString()+"/inventory/adjust", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
	if service.userID == nil || *service.userID != authenticatedUserID {
		t.Fatalf("expected audit user %s, got %v", authenticatedUserID, service.userID)
	}
}
