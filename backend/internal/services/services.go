package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

// Services holds all service interfaces
type Services struct {
	Auth      AuthService
	User      UserService
	Category  CategoryService
	Inbound   InboundService
	Outbound  OutboundService
	Inventory InventoryService
	Report    ReportService
}

// NewServices creates a new services instance
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth:      NewAuthService(repos.User),
		User:      NewUserService(repos.User),
		Category:  NewCategoryService(repos.Category, repos.Inventory),
		Inbound:   NewInboundService(repos.Inbound, repos.Inventory),
		Outbound:  NewOutboundService(repos.Outbound, repos.Inventory),
		Inventory: NewInventoryService(repos.Inventory, repos.Category),
		Report:    NewReportService(repos),
	}
}

// AuthService interface
type AuthService interface {
	Login(username, password string) (*models.LoginResponse, error)
	ValidateToken(tokenString string) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
}

// UserService interface
type UserService interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// CategoryService interface
type CategoryService interface {
	Create(category *models.BatteryCategory) error
	GetByID(id uint) (*models.BatteryCategory, error)
	GetAll() ([]models.BatteryCategory, error)
	Update(category *models.BatteryCategory) error
	Delete(id uint) error
}

// InboundService interface
type InboundService interface {
	Create(req *models.CreateInboundOrderRequest, createdBy uint) (*models.InboundOrder, error)
	GetByID(id uint) (*models.InboundOrder, error)
	GetAll(limit, offset int) ([]models.InboundOrder, int64, error)
	Update(order *models.InboundOrder) error
	Delete(id uint) error
}

// OutboundService interface
type OutboundService interface {
	Create(req *models.CreateOutboundOrderRequest, createdBy uint) (*models.OutboundOrder, error)
	GetByID(id uint) (*models.OutboundOrder, error)
	GetAll(limit, offset int) ([]models.OutboundOrder, int64, error)
	Update(order *models.OutboundOrder) error
	Delete(id uint) error
}

// InventoryService interface
type InventoryService interface {
	GetByCategoryID(categoryID uint) (*models.Inventory, error)
	GetAll() ([]models.Inventory, error)
}

// ReportService interface
type ReportService interface {
	GetSummary(startDate, endDate string) (*models.ReportSummary, error)
}
