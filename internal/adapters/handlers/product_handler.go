package handlers

import (
	"net/http"
	"strconv"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles HTTP requests for products.
type ProductHandler struct {
	service ports.ProductService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetAll returns all products for a store.
func (h *ProductHandler) GetAll(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	products, err := h.service.GetAllProducts(uint(storeID))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, products)
}

// GetByID returns a single product.
func (h *ProductHandler) GetByID(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	product, err := h.service.GetProduct(uint(storeID), uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Product not found")
		return
	}
	pkg.Success(c, http.StatusOK, product)
}

// Create adds a new product.
func (h *ProductHandler) Create(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	product.StoreID = uint(storeID)

	if err := h.service.CreateProduct(&product); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusCreated, product)
}

// Update modifies an existing product.
func (h *ProductHandler) Update(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	product, err := h.service.GetProduct(uint(storeID), uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	if err := c.ShouldBindJSON(product); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	product.StoreID = uint(storeID)
	product.ID = uint(id)

	if err := h.service.UpdateProduct(product); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, product)
}

// Delete soft-deletes a product.
func (h *ProductHandler) Delete(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.service.DeleteProduct(uint(storeID), uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message": "Product deleted"})
}
