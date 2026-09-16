package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/anoulack007/core-pos/config"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/dto"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service     ports.AuthService
	minioClient *minio.Client
	cfg         *config.Config
}

func NewAuthHandler(service ports.AuthService, mc *minio.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{service: service, minioClient: mc, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	isMultipart := strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data")
	if isMultipart {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 11<<20)
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			pkg.Error(c, http.StatusBadRequest, "file too large or invalid multipart form")
			return
		}
		req = dto.RegisterRequest{
			StoreID:  c.PostForm("store_id"),
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
			FullName: c.PostForm("full_name"),
			Email:    c.PostForm("email"),
			Phone:    c.PostForm("phone"),
		}
	} else if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.StoreID == "" || req.Username == "" || req.Password == "" {
		pkg.Error(c, http.StatusBadRequest, "store_id, username and password are required")
		return
	}
	if len(req.Password) < 8 {
		pkg.Error(c, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	storeID, err := uuid.Parse(req.StoreID)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid store ID")
		return
	}

	var avatarURL string

	if isMultipart {
		file, fileErr := c.FormFile("avatar")
		if fileErr == nil {
			if file.Size > 10<<20 {
				pkg.Error(c, http.StatusBadRequest, "avatar must not exceed 10 MB")
				return
			}
			filename := fmt.Sprintf("avatars/%s%s", uuid.New().String(), filepath.Ext(file.Filename))
			src, err := file.Open()
			if err != nil {
				pkg.Error(c, http.StatusBadRequest, "invalid avatar file")
				return
			}
			defer src.Close()

			_, err = h.minioClient.PutObject(context.Background(), h.cfg.MinioBucket, filename, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
			if err != nil {
				pkg.Error(c, http.StatusInternalServerError, "failed to upload avatar: "+err.Error())
				return
			}
			avatarURL = "/" + h.cfg.MinioBucket + "/" + filename
		} else if !errors.Is(fileErr, http.ErrMissingFile) {
			pkg.Error(c, http.StatusBadRequest, "invalid avatar file")
			return
		}
	}

	user := domain.User{
		StoreID:   storeID,
		Username:  req.Username,
		Role:      domain.RoleOwner,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		AvatarURL: avatarURL,
	}

	if err := h.service.Register(&user, req.Password); err != nil {
		if errors.Is(err, domain.ErrStoreAlreadyInitialized) {
			pkg.Error(c, http.StatusConflict, "store registration is already complete")
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			pkg.Error(c, http.StatusNotFound, "store not found")
			return
		}
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.Success(c, http.StatusCreated, gin.H{
		"message":    "user created successfully",
		"avatar_url": user.AvatarURL,
	})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// fmt.Println(req.Username, req.Password)
	accessToken, refreshToken, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		pkg.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	accessToken, newRefreshToken, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		pkg.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.service.Logout(); err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"message": "logged out successfully"})
}
