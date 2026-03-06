package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/anoulack007/core-pos/config"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadHandler struct {
	minioClient *minio.Client
	cfg			*config.Config
}

func NewUploadHandler(mc *minio.Client, cfg *config.Config) *UploadHandler {
	return &UploadHandler{minioClient: mc, cfg: cfg}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	folder := c.DefaultPostForm("folder", "other")
	filename := fmt.Sprintf("%s/%s%s",folder, uuid.New().String(),filepath.Ext(file.Filename))
	src, _ := file.Open()
	defer src.Close()

	_, err := h.minioClient.PutObject(context.Background(), h.cfg.MinioBucket, filename,src,
	file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(c, http.StatusOK, gin.H{"url":"/" + h.cfg.MinioBucket + "/" + filename})
}