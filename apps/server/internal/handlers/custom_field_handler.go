package handlers

import (
	"context"
	"net/http"
	"strconv"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"

	"github.com/gin-gonic/gin"
)

type CustomFieldHandler struct {
	service CustomFieldService
}

type CustomFieldService interface {
	List(ctx context.Context, resource string, activeOnly bool) ([]models.CustomField, error)
	Get(ctx context.Context, id uint) (*models.CustomField, error)
	Create(ctx context.Context, req *services.CustomFieldCreateRequest) (*models.CustomField, error)
	Update(ctx context.Context, id uint, req *services.CustomFieldUpdateRequest) (*models.CustomField, error)
	Delete(ctx context.Context, id uint) error
}

func NewCustomFieldHandler(service CustomFieldService) *CustomFieldHandler {
	return &CustomFieldHandler{service: service}
}

func (h *CustomFieldHandler) List(c *gin.Context) {
	resource := c.DefaultQuery("resource", "ticket")
	activeOnly := c.DefaultQuery("active", "true") != "false"
	fields, err := h.service.List(c.Request.Context(), resource, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list custom fields", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, fields)
}

func (h *CustomFieldHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id", Message: err.Error()})
		return
	}
	field, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Custom field not found", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, field)
}

func (h *CustomFieldHandler) Create(c *gin.Context) {
	var req services.CustomFieldCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	field, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to create custom field", Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, field)
}

func (h *CustomFieldHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id", Message: err.Error()})
		return
	}
	var req services.CustomFieldUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	field, err := h.service.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		status := http.StatusBadRequest
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{Error: "Failed to update custom field", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, field)
}

func (h *CustomFieldHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid id", Message: err.Error()})
		return
	}
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		status := http.StatusBadRequest
		if isNotFoundError(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, ErrorResponse{Error: "Failed to delete custom field", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, SuccessResponse{Message: "deleted"})
}

func RegisterCustomFieldRoutes(r *gin.RouterGroup, handler *CustomFieldHandler) {
	fields := r.Group("/custom-fields")
	{
		fields.GET("", handler.List)
		fields.GET("/:id", handler.Get)
		fields.POST("", handler.Create)
		fields.PUT("/:id", handler.Update)
		fields.DELETE("/:id", handler.Delete)
	}
}
