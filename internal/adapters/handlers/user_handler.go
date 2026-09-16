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
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	users, err := h.service.GetStaff(storeID)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, "failed to load staff")
		return
	}
	pkg.Success(c, http.StatusOK, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	storeID, err := uuid.Parse(c.Param("storeId"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
		return
	}

	roleValue, exists := c.Get(middleware.ContextRoleKey)
	actorRole, ok := roleValue.(domain.UserRole)
	if !exists || !ok {
		pkg.Error(c, http.StatusForbidden, "role access denied")
		return
	}

	var req dto.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := domain.User{
		StoreID:  storeID,
		Username: req.Username,
		Role:     domain.UserRole(req.Role),
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	if err := h.service.CreateStaff(actorRole, &user, req.Password); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbiddenRoleAssignment):
			pkg.Error(c, http.StatusForbidden, err.Error())
		case errors.Is(err, domain.ErrInvalidStaff):
			pkg.Error(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrUsernameAlreadyExists):
			pkg.Error(c, http.StatusConflict, err.Error())
		default:
			pkg.Error(c, http.StatusInternalServerError, "failed to create staff")
		}
		return
	}

	pkg.Success(c, http.StatusCreated, user)
}
