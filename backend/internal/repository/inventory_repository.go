package repository

import (
	"battery-erp-backend/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByCategoryID(categoryID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("category_id = ?", categoryID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) GetAll() ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.Find(&inventories).Error
	return inventories, err
}

func (r *inventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

func (r *inventoryRepository) Update(inventory *models.Inventory) error {
	return r.db.Save(inventory).Error
}

func (r *inventoryRepository) UpdateWeight(categoryID uint, weightChange float64, isInbound bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory

		// Find or create inventory record
		err := tx.Where("category_id = ?", categoryID).First(&inventory).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new inventory record
				now := time.Now()
				inventory = models.Inventory{
					CategoryID:      categoryID,
					CurrentWeightKg: 0,
					CreatedAt:       now,
					UpdatedAt:       now,
				}
				if err := tx.Create(&inventory).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Update weight
		if isInbound {
			inventory.CurrentWeightKg += weightChange
			now := time.Now()
			inventory.LastInboundAt = &now
		} else {
			// Check for overselling
			if inventory.CurrentWeightKg < weightChange {
				return errors.New(fmt.Sprintf("insufficient inventory: current weight is %.3f kg, requested %.3f kg",
					inventory.CurrentWeightKg, weightChange))
			}
			inventory.CurrentWeightKg -= weightChange
			now := time.Now()
			inventory.LastOutboundAt = &now
		}

		inventory.UpdatedAt = time.Now()
		return tx.Save(&inventory).Error
	})
}
