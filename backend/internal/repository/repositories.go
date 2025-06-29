package repository

import (
	"battery-erp-backend/internal/models"

	"gorm.io/gorm"
)

// Repositories holds all repository interfaces
type Repositories struct {
	User      UserRepository
	Category  CategoryRepository
	Inbound   InboundRepository
	Outbound  OutboundRepository
	Inventory InventoryRepository
}

// NewRepositories creates a new repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db),
		Category:  NewCategoryRepository(db),
		Inbound:   NewInboundRepository(db),
		Outbound:  NewOutboundRepository(db),
		Inventory: NewInventoryRepository(db),
	}
}

// UserRepository interface
type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// CategoryRepository interface
type CategoryRepository interface {
	Create(category *models.BatteryCategory) error
	GetByID(id uint) (*models.BatteryCategory, error)
	GetAll() ([]models.BatteryCategory, error)
	Update(category *models.BatteryCategory) error
	Delete(id uint) error
}

// InboundRepository interface
type InboundRepository interface {
	Create(order *models.InboundOrder) error
	CreateItem(item *models.InboundOrderItem) error
	GetByID(id uint) (*models.InboundOrder, error)
	GetAll(limit, offset int) ([]models.InboundOrder, int64, error)
	Update(order *models.InboundOrder) error
	Delete(id uint) error
	GenerateOrderNo() (string, error)
}

// OutboundRepository interface
type OutboundRepository interface {
	Create(order *models.OutboundOrder) error
	CreateItem(item *models.OutboundOrderItem) error
	GetByID(id uint) (*models.OutboundOrder, error)
	GetAll(limit, offset int) ([]models.OutboundOrder, int64, error)
	Update(order *models.OutboundOrder) error
	Delete(id uint) error
	GenerateOrderNo() (string, error)
}

// InventoryRepository interface
type InventoryRepository interface {
	GetByCategoryID(categoryID uint) (*models.Inventory, error)
	GetAll() ([]models.Inventory, error)
	Create(inventory *models.Inventory) error
	Update(inventory *models.Inventory) error
	UpdateWeight(categoryID uint, weightChange float64, isInbound bool) error
}
