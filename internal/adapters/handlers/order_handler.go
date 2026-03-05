package handlers

import (
	"net/http"
	"strconv"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	service ports.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(service ports.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// GetAll returns all orders for a store.
func (h *OrderHandler) GetAll(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	orders, err := h.service.GetAllOrders(uint(storeID))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, orders)
}

// GetByID returns a single order with items and payments.
func (h *OrderHandler) GetByID(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	order, err := h.service.GetOrder(uint(storeID), uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Order not found")
		return
	}
	pkg.Success(c, http.StatusOK, order)
}

// Create creates a new order.
func (h *OrderHandler) Create(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	order.StoreID = uint(storeID)

	if err := h.service.CreateOrder(&order); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusCreated, order)
}

// Void voids an existing order.
func (h *OrderHandler) Void(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.service.VoidOrder(uint(storeID), uint(id)); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message": "Order voided"})
}
