package services

import (
	"battery-erp-backend/internal/repository"

	"gorm.io/gorm"
)

// Services holds all service instances (no interfaces)
type Services struct {
	UserService      *UserService
	CategoryService  *CategoryService
	InboundService   *InboundService
	OutboundService  *OutboundService
	InventoryService *InventoryService
	SellerService    *SellerService
	ReportService    *ReportService
	Auth             *AuthService
	DB               *gorm.DB
}

// NewServices creates a new services instance
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		UserService:      NewUserService(repos.UserRepo),
		CategoryService:  NewCategoryService(repos.CategoryRepo, repos.InventoryRepo),
		InboundService:   NewInboundService(repos.InboundRepo, repos.InventoryRepo),
		OutboundService:  NewOutboundService(repos.OutboundRepo, repos.InventoryRepo),
		InventoryService: NewInventoryService(repos.InventoryRepo, repos.CategoryRepo),
		SellerService:    NewSellerService(repos.SellerRepo),
		ReportService:    NewReportService(repos),
		Auth:             NewAuthService(repos.UserRepo),
		DB:               repos.DB,
	}
}
