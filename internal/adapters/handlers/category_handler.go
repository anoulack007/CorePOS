package handlers

import (
	"net/http"
	"strconv"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
)

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	service ports.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(service ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll returns all categories for a store.
func (h *CategoryHandler) GetAll(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	categories, err := h.service.GetAllCategories(uint(storeID))
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, categories)
}

// GetByID returns a single category.
func (h *CategoryHandler) GetByID(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	category, err := h.service.GetCategory(uint(storeID), uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Category not found")
		return
	}
	pkg.Success(c, http.StatusOK, category)
}

// Create adds a new category.
func (h *CategoryHandler) Create(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)

	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	category.StoreID = uint(storeID)

	if err := h.service.CreateCategory(&category); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusCreated, category)
}

// Update modifies an existing category.
func (h *CategoryHandler) Update(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	category, err := h.service.GetCategory(uint(storeID), uint(id))
	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	if err := c.ShouldBindJSON(category); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	category.StoreID = uint(storeID)
	category.ID = uint(id)

	if err := h.service.UpdateCategory(category); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, category)
}

// Delete removes a category.
func (h *CategoryHandler) Delete(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("storeId"), 10, 32)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.service.DeleteCategory(uint(storeID), uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message": "Category deleted"})
}
