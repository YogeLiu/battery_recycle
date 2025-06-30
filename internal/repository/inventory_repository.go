package repository

import (
	"battery-erp-backend/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// InventoryRepository 库存数据仓库 (不再是接口实现)
type InventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository 创建库存仓库实例
func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

// GetByCategoryID 根据分类ID获取库存
func (r *InventoryRepository) GetByCategoryID(categoryID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	err := r.db.Where("category_id = ?", categoryID).First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

// GetAll 获取所有库存
func (r *InventoryRepository) GetAll() ([]models.Inventory, error) {
	var inventories []models.Inventory
	err := r.db.Find(&inventories).Error
	return inventories, err
}

// Create 创建库存记录
func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	return r.db.Create(inventory).Error
}

// UpdateCurrentWeight 显式更新当前重量
func (r *InventoryRepository) UpdateCurrentWeight(categoryID uint, weight float64) error {
	return r.db.Model(&models.Inventory{}).Where("category_id = ?", categoryID).Update("current_weight_kg", weight).Error
}

// UpdateLastInboundAt 显式更新最后入库时间
func (r *InventoryRepository) UpdateLastInboundAt(categoryID uint, lastInboundAt time.Time) error {
	return r.db.Model(&models.Inventory{}).Where("category_id = ?", categoryID).Update("last_inbound_at", lastInboundAt).Error
}

// UpdateLastOutboundAt 显式更新最后出库时间
func (r *InventoryRepository) UpdateLastOutboundAt(categoryID uint, lastOutboundAt time.Time) error {
	return r.db.Model(&models.Inventory{}).Where("category_id = ?", categoryID).Update("last_outbound_at", lastOutboundAt).Error
}

// UpdateFields 显式更新指定字段
func (r *InventoryRepository) UpdateFields(categoryID uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Inventory{}).Where("category_id = ?", categoryID).Updates(updates).Error
}

// UpdateWeight 显式更新库存重量 (事务)
func (r *InventoryRepository) UpdateWeight(categoryID uint, weightChange float64, isInbound bool) error {
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

		// 显式更新相关字段
		now := time.Now()
		updates := map[string]interface{}{
			"updated_at": now,
		}

		if isInbound {
			newWeight := inventory.CurrentWeightKg + weightChange
			updates["current_weight_kg"] = newWeight
			updates["last_inbound_at"] = now
		} else {
			// Check for overselling
			if inventory.CurrentWeightKg < weightChange {
				return errors.New(fmt.Sprintf("insufficient inventory: current weight is %.3f kg, requested %.3f kg",
					inventory.CurrentWeightKg, weightChange))
			}
			newWeight := inventory.CurrentWeightKg - weightChange
			updates["current_weight_kg"] = newWeight
			updates["last_outbound_at"] = now
		}

		return tx.Model(&models.Inventory{}).Where("category_id = ?", categoryID).Updates(updates).Error
	})
}
