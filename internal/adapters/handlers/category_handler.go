package handlers

import (
	"net/http"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	service ports.CategoryService
}

func NewCategoryHandler(service ports.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetAll(c *gin.Context){
	storeID, err := uuid.Parse(c.Param("storeId"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	categories, err := h.service.GetAllCategories(storeID)

	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest,"Invalid category ID")
		return
	}

	category, err := h.service.GetCategory(storeID,id)
	if err != nil {
		pkg.Error(c, http.StatusNotFound,"Category not found")
		return
	}

	pkg.Success(c, http.StatusOK,category)
}


func (h *CategoryHandler) Create(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		pkg.Error(c, http.StatusBadRequest,err.Error())
		return
	}

	category.StoreID = storeID

	if err := h.service.CreateCategory(&category); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil{
		pkg.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}


	category, err := h.service.GetCategory(storeID,id)

	if err != nil {
		pkg.Error(c, http.StatusNotFound, "Category not found")
		return
	}

	if err := c.ShouldBindJSON(category); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category.ID = id
	category.StoreID = storeID

	if err := h.service.UpdateCategory(category); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c,http.StatusOK, category)
}


func (h *CategoryHandler) Delete(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.DeleteCategory(storeID, id); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusOK, gin.H{"message": "Category deleted"})
}