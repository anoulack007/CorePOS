package handlers

import (
	"errors"
	"net/http"

	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service ports.InventoryService
}

func NewInventoryHandler(service ports.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

type AdjustStockRequest struct {
	ProductID       uuid.UUID `json:"product_id" binding:"required"`
	MovementType    string    `json:"movement_type" binding:"required"`
	QuantityChanged int       `json:"quantity_changed" binding:"required"`
	Notes           string    `json:"notes"`
	EvidenceURL     string    `json:"evidence_url"`
}

func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	var req AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		pkg.Error(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	authenticatedUserID, ok := userID.(uuid.UUID)
	if !ok {
		pkg.Error(c, http.StatusUnauthorized, "invalid authenticated user")
		return
	}

	if err := h.service.AdjustStock(storeID, req.ProductID, &authenticatedUserID, req.MovementType, req.QuantityChanged, req.Notes, req.EvidenceURL); err != nil {
		if errors.Is(err, domain.ErrInvalidInventoryMovement) {
			pkg.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			pkg.Error(c, http.StatusConflict, err.Error())
			return
		}
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusOK, gin.H{"message": "Stock adjusted successfully"})
}

func (h *InventoryHandler) GetHistory(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	var productIDPtr *uuid.UUID
	productIDStr := c.Query("productId")
	if productIDStr != "" {
		pid, err := uuid.Parse(productIDStr)
		if err == nil {
			productIDPtr = &pid
		} else {
			pkg.Error(c, http.StatusBadRequest, "Invalid product ID format")
			return
		}
	}

	history, err := h.service.GetStockHistory(storeID, productIDPtr)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusOK, history)
}
