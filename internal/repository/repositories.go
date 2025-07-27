package repository

import (
	"battery-erp-backend/internal/models"

	"gorm.io/gorm"
)

// Repositories holds all repository instances (no interfaces)
type Repositories struct {
	UserRepo      *UserRepository
	CategoryRepo  *CategoryRepository
	InboundRepo   *InboundRepository
	OutboundRepo  *OutboundRepository
	InventoryRepo *InventoryRepository
	SellerRepo    *SellerRepository
	DB            *gorm.DB
}

// NewRepositories creates a new repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:      NewUserRepository(db),
		CategoryRepo:  NewCategoryRepository(db),
		InboundRepo:   NewInboundRepository(db),
		OutboundRepo:  NewOutboundRepository(db),
		InventoryRepo: NewInventoryRepository(db),
		SellerRepo:    NewSellerRepository(db),
		DB:            db,
	}
}

// AutoMigrate 自动迁移数据库表结构
func (r *Repositories) AutoMigrate() error {
	return r.DB.AutoMigrate(
		&models.User{},
		&models.BatteryCategory{},
		&models.InboundOrder{},
		&models.OutboundOrder{},
		&models.Inventory{},
		&models.Seller{},
	)
}
