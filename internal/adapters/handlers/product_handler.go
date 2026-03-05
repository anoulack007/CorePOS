package handlers

import (
	"net/http"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler{
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(c *gin.Context){
	storeID,err := uuid.Parse(c.Param("storeId"))
	if err != nil{
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}
	products,err := h.service.GetAllProducts(storeID)
	if err != nil{
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(c *gin.Context){
	storeID, _ := uuid.Parse(c.Param("storeId"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil{
		pkg.Error(c, http.StatusBadRequest, "Invalid product ID")
		return
	}
	product, err:= h.service.GetProduct(storeID, id)
	if err != nil{
		pkg.Error(c, http.StatusNotFound, "Product not found")
		return
	}
	pkg.Success(c, http.StatusOK, product)
}

func (h *ProductHandler) Create(c *gin.Context){
	storeID,_ := uuid.Parse(c.Param("storeId"))
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	product.StoreID = storeID
	if err := h.service.CreateProduct(&product); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusCreated, product)
}

func (h *ProductHandler) Update(c *gin.Context){
	storeID, _ := uuid.Parse(c.Param("storeId"))
	id,_ := uuid.Parse(c.Param("id"))
	product, err := h.service.GetProduct(storeID, id)
	if err != nil{
		pkg.Error(c, http.StatusNotFound,"Product not found")
		return
	}

	if err := c.ShouldBindJSON(product); err != nil{
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	product.StoreID = storeID

	product.ID = id

	if err:=h.service.UpdateProduct(product); err !=nil{
		pkg.Error(c,http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c,http.StatusOK,product)
}

func (h *ProductHandler) Delete(c *gin.Context){
	storeID,_ := uuid.Parse(c.Param("storeId"))
	id,_ := uuid.Parse(c.Param("id"))
	if err := h.service.DeleteProduct(storeID,id); err != nil{
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message":"Product deleted"})
}