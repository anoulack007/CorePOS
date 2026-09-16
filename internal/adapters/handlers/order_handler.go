package handlers

import (
	"errors"
	"net/http"

	"github.com/anoulack007/core-pos/internal/adapters/middleware"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/dto"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderHandler struct {
	service ports.OrderService
}

func NewOrderHandler(service ports.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}
	orders, err := h.service.GetAllOrders(storeID)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, "failed to load orders")
		return
	}
	pkg.Success(c, http.StatusOK, orders)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	storeID, orderID, ok := parseOrderIDs(c)
	if !ok {
		return
	}
	order, err := h.service.GetOrder(storeID, orderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pkg.Error(c, http.StatusNotFound, "Order not found")
		return
	}
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, "failed to load order")
		return
	}
	pkg.Success(c, http.StatusOK, order)
}

func (h *OrderHandler) Create(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}
	userValue, exists := c.Get(middleware.ContextUserIDKey)
	userID, ok := userValue.(uuid.UUID)
	if !exists || !ok {
		pkg.Error(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	lines := make([]ports.OrderLineInput, 0, len(req.Items))
	for _, item := range req.Items {
		lines = append(lines, ports.OrderLineInput{ProductID: item.ProductID, Quantity: item.Quantity})
	}

	order, err := h.service.CreateOrder(storeID, userID, lines)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOrder):
			pkg.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrProductNotFound):
			pkg.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInsufficientStock):
			pkg.Error(c, http.StatusConflict, err.Error())
		default:
			pkg.Error(c, http.StatusInternalServerError, "failed to create order")
		}
		return
	}
	pkg.Success(c, http.StatusCreated, order)
}

func (h *OrderHandler) Void(c *gin.Context) {
	storeID, orderID, ok := parseOrderIDs(c)
	if !ok {
		return
	}
	userValue, exists := c.Get(middleware.ContextUserIDKey)
	userID, ok := userValue.(uuid.UUID)
	if !exists || !ok {
		pkg.Error(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	if err := h.service.VoidOrder(storeID, orderID, userID); err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			pkg.Error(c, http.StatusNotFound, "Order not found")
		case errors.Is(err, domain.ErrOrderNotVoidable), errors.Is(err, domain.ErrPaidOrderCannotVoid):
			pkg.Error(c, http.StatusConflict, err.Error())
		default:
			pkg.Error(c, http.StatusInternalServerError, "failed to void order")
		}
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message": "Order voided"})
}

func parseOrderIDs(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return uuid.Nil, uuid.Nil, false
	}
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid order ID")
		return uuid.Nil, uuid.Nil, false
	}
	return storeID, orderID, true
}
