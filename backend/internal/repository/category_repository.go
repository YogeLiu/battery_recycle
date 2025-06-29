package repository

import (
	"battery-erp-backend/internal/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.BatteryCategory) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*models.BatteryCategory, error) {
	var category models.BatteryCategory
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll() ([]models.BatteryCategory, error) {
	var categories []models.BatteryCategory
	err := r.db.Where("is_active = ?", true).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *models.BatteryCategory) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Model(&models.BatteryCategory{}).Where("id = ?", id).Update("is_active", false).Error
}