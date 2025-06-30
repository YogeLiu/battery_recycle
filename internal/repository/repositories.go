package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository instances (no interfaces)
type Repositories struct {
	User      *UserRepository
	Category  *CategoryRepository
	Inbound   *InboundRepository
	Outbound  *OutboundRepository
	Inventory *InventoryRepository
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
