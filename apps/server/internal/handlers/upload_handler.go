package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"servify/apps/server/internal/platform/storage"

	"github.com/gin-gonic/gin"
)

// FileUploadHandler handles file upload endpoints using the storage.Provider abstraction.
type FileUploadHandler struct {
	provider storage.Provider
	maxSize  int64 // bytes
}

// NewFileUploadHandler creates a new FileUploadHandler backed by a storage.Provider.
func NewFileUploadHandler(provider storage.Provider, maxSize int64) *FileUploadHandler {
	return &FileUploadHandler{provider: provider, maxSize: maxSize}
}

// Upload handles file upload.
func (h *FileUploadHandler) Upload(c *gin.Context) {
	if h.maxSize > 0 {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxSize)
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件", "message": err.Error()})
		return
	}
	defer file.Close()

	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".txt": true, ".csv": true, ".zip": true, ".mp3": true, ".mp4": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件类型", "message": fmt.Sprintf("文件类型 %s 不被允许", ext)})
		return
	}

	// Build storage key with date partition
	dateDir := time.Now().Format("2006/01/02")
	key := dateDir + "/" + fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))

	info, err := h.provider.Save(key, file, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "上传成功",
		"filename": header.Filename,
		"url":      info.URL,
		"size":     info.Size,
	})
}
