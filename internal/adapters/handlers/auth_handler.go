package handlers

import (
	// "fmt"
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/anoulack007/core-pos/config"
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/dto"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "file too large or invalid multipart form")
		return
	}

	storeIDStr := c.PostForm("store_id")
	username := c.PostForm("username")
	password := c.PostForm("password")
	role := c.PostForm("role")
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")

	if storeIDStr == "" || username == "" || password == "" {
		pkg.Error(c, http.StatusBadRequest, "store_id, username and password are required")
		return
	}

	storeID, err := uuid.Parse(storeIDStr)
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, "invalid store ID")
		return
	}

	var avatarURL string

	file, err := c.FormFile("avatar")

	if err == nil {
		filename := fmt.Sprintf("avatars/%s%s", uuid.New().String(), filepath.Ext(file.Filename))
		src, _ := file.Open()
		defer src.Close()

		_, err = h.minioClient.PutObject(context.Background(), h.cfg.MinioBucket, filename, src, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})

		if err != nil {
			pkg.Error(c, http.StatusInternalServerError, "failed to upload avatar: "+err.Error())
			return
		}
		avatarURL = "/" + h.cfg.MinioBucket + "/" + filename
	}

	user := domain.User{
		StoreID:   storeID,
		Username:  username,
		Role:      domain.NoralizeUserRole(role),
		FullName:  fullName,
		Email:     email,
		Phone:     phone,
		AvatarURL: avatarURL,
	}

	if err := h.service.Register(&user, password); err != nil {
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
