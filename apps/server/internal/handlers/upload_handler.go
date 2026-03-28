package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// FileUploadHandler handles file upload endpoints.
type FileUploadHandler struct {
	uploadDir string
	maxSize   int64 // bytes
}

// NewFileUploadHandler creates a new FileUploadHandler.
func NewFileUploadHandler(uploadDir string, maxSize int64) *FileUploadHandler {
	return &FileUploadHandler{uploadDir: uploadDir, maxSize: maxSize}
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

	// Generate unique filename
	dateDir := time.Now().Format("2006/01/02")
	saveDir := filepath.Join(h.uploadDir, dateDir)
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	savePath := filepath.Join(saveDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "写入文件失败"})
		return
	}

	url := "/uploads/" + dateDir + "/" + filename
	c.JSON(http.StatusCreated, gin.H{
		"message":  "上传成功",
		"filename": header.Filename,
		"url":      url,
		"size":     header.Size,
	})
}
