package repository

import (
	"battery-erp-backend/internal/models"

	"gorm.io/gorm"
)

// CategoryRepository 电池类别数据仓库 (不再是接口实现)
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建类别仓库实例
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建电池类别
func (r *CategoryRepository) Create(category *models.BatteryCategory) error {
	return r.db.Create(category).Error
}

// GetByID 根据ID获取电池类别
func (r *CategoryRepository) GetByID(id uint) (*models.BatteryCategory, error) {
	var category models.BatteryCategory
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll 获取所有活跃的电池类别
func (r *CategoryRepository) GetAll() ([]models.BatteryCategory, error) {
	var categories []models.BatteryCategory
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

// UpdateName 显式更新类别名称
func (r *CategoryRepository) UpdateName(id uint, name string) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Update("name", name).Error
}

// UpdateDescription 显式更新类别描述
func (r *CategoryRepository) UpdateDescription(id uint, description string) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Update("description", description).Error
}

// UpdateUnitPrice 显式更新单价
func (r *CategoryRepository) UpdateUnitPrice(id uint, unitPrice float64) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Update("unit_price", unitPrice).Error
}

// UpdateFields 显式更新指定字段
func (r *CategoryRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除类别 (设置为非活跃状态)
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Update("is_active", false).Error
}
