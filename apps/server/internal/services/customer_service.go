package services

import (
	"context"
	"strings"

	"servify/apps/server/internal/models"
	customerapp "servify/apps/server/internal/modules/customer/application"
	customerinfra "servify/apps/server/internal/modules/customer/infra"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CustomerService 客户管理服务兼容层。
type CustomerService struct {
	db     *gorm.DB
	logger *logrus.Logger
	module *customerapp.Service
}

// NewCustomerService 创建客户服务。
func NewCustomerService(db *gorm.DB, logger *logrus.Logger) *CustomerService {
	if logger == nil {
		logger = logrus.New()
	}
	return &CustomerService{
		db:     db,
		logger: logger,
		module: customerapp.NewService(customerinfra.NewGormRepository(db)),
	}
}

// CustomerCreateRequest 创建客户请求。
type CustomerCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
	Industry string `json:"industry"`
	Source   string `json:"source"`
	Tags     string `json:"tags"`
	Notes    string `json:"notes"`
	Priority string `json:"priority"`
}

// CustomerUpdateRequest 更新客户请求。
type CustomerUpdateRequest struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Company  *string `json:"company"`
	Industry *string `json:"industry"`
	Source   *string `json:"source"`
	Tags     *string `json:"tags"`
	Notes    *string `json:"notes"`
	Priority *string `json:"priority"`
	Status   *string `json:"status"`
}

// CustomerListRequest 客户列表请求。
type CustomerListRequest struct {
	Page      int      `form:"page,default=1"`
	PageSize  int      `form:"page_size,default=20"`
	Search    string   `form:"search"`
	Industry  []string `form:"industry"`
	Source    []string `form:"source"`
	Priority  []string `form:"priority"`
	Status    []string `form:"status"`
	Tags      string   `form:"tags"`
	SortBy    string   `form:"sort_by,default=created_at"`
	SortOrder string   `form:"sort_order,default=desc"`
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *CustomerCreateRequest) (*models.User, error) {
	return s.module.CreateCustomer(ctx, customerapp.CreateCustomerCommand{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Company:  req.Company,
		Industry: req.Industry,
		Source:   req.Source,
		Tags:     splitTags(req.Tags),
		Notes:    req.Notes,
		Priority: req.Priority,
	})
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error) {
	return s.module.GetCustomerByID(ctx, customerID)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, customerID uint, req *CustomerUpdateRequest) (*models.User, error) {
	cmd := customerapp.UpdateCustomerCommand{
		Name:     req.Name,
		Phone:    req.Phone,
		Company:  req.Company,
		Industry: req.Industry,
		Source:   req.Source,
		Notes:    req.Notes,
		Priority: req.Priority,
		Status:   req.Status,
	}
	if req.Tags != nil {
		tags := splitTags(*req.Tags)
		cmd.Tags = &tags
	}
	return s.module.UpdateCustomer(ctx, customerID, cmd)
}

func (s *CustomerService) ListCustomers(ctx context.Context, req *CustomerListRequest) ([]CustomerInfo, int64, error) {
	items, total, err := s.module.ListCustomers(ctx, customerapp.ListCustomersQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Search:    req.Search,
		Industry:  req.Industry,
		Source:    req.Source,
		Priority:  req.Priority,
		Status:    req.Status,
		Tags:      splitTags(req.Tags),
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, 0, err
	}
	out := make([]CustomerInfo, 0, len(items))
	for _, item := range items {
		out = append(out, CustomerInfo(item))
	}
	return out, total, nil
}

func (s *CustomerService) GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*CustomerActivity, error) {
	activity, err := s.module.GetCustomerActivity(ctx, customerID, limit)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *CustomerService) AddCustomerNote(ctx context.Context, customerID uint, note string, userID uint) error {
	return s.module.AddNote(ctx, customerID, note, userID)
}

func (s *CustomerService) UpdateCustomerTags(ctx context.Context, customerID uint, tags []string) error {
	return s.module.UpdateTags(ctx, customerID, tags)
}

func (s *CustomerService) GetCustomerStats(ctx context.Context) (*CustomerStats, error) {
	stats, err := s.module.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// CustomerInfo 客户信息（用于列表显示）。
type CustomerInfo = customerapp.CustomerInfoDTO

// CustomerActivity 客户活动记录。
type CustomerActivity = customerapp.CustomerActivityDTO

// CustomerStats 客户统计信息。
type CustomerStats = customerapp.CustomerStatsDTO

func splitTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
