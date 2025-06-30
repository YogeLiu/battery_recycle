package services

import (
	"battery-erp-backend/internal/repository"
)

// Services holds all service instances (no interfaces)
type Services struct {
	Auth      *AuthService
	User      *UserService
	Category  *CategoryService
	Inbound   *InboundService
	Outbound  *OutboundService
	Inventory *InventoryService
	Report    *ReportService
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
