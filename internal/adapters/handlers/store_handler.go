package handlers

import (
	"net/http"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StoreHandler struct {
	db *gorm.DB
}

func NewStoreHandler(db *gorm.DB) *StoreHandler {
	return &StoreHandler{db: db}
}

func (h *StoreHandler) Create(c *gin.Context) {
	var store domain.Store
	if err := c.ShouldBindJSON(&store); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.db.Create(&store).Error; err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusCreated, store)
}

func (h *StoreHandler) GetAll(c *gin.Context) {
	var stores []domain.Store
	if err := h.db.Find(&stores).Error; err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, stores)
}
