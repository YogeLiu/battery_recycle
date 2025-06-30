package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

// CategoryService 类别服务 (不再使用接口)
type CategoryService struct {
	categoryRepo  *repository.CategoryRepository
	inventoryRepo *repository.InventoryRepository
}

// NewCategoryService 创建类别服务实例
func NewCategoryService(categoryRepo *repository.CategoryRepository, inventoryRepo *repository.InventoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo:  categoryRepo,
		inventoryRepo: inventoryRepo,
	}
}

// Create 创建分类
func (s *CategoryService) Create(category *models.BatteryCategory) error {
	// 创建分类
	if err := s.categoryRepo.Create(category); err != nil {
		return err
	}

	// 为新分类初始化库存记录
	inventory := &models.Inventory{
		CategoryID:      category.ID,
		CurrentWeightKg: 0,
	}
	return s.inventoryRepo.Create(inventory)
}

// GetByID 根据ID获取分类
func (s *CategoryService) GetByID(id uint) (*models.BatteryCategory, error) {
	return s.categoryRepo.GetByID(id)
}

// GetAll 获取所有分类
func (s *CategoryService) GetAll() ([]models.BatteryCategory, error) {
	return s.categoryRepo.GetAll()
}

// UpdateName 显式更新分类名称
func (s *CategoryService) UpdateName(id uint, name string) error {
	return s.categoryRepo.UpdateName(id, name)
}

// UpdateDescription 显式更新分类描述
func (s *CategoryService) UpdateDescription(id uint, description string) error {
	return s.categoryRepo.UpdateDescription(id, description)
}

// UpdateUnitPrice 显式更新单价
func (s *CategoryService) UpdateUnitPrice(id uint, unitPrice float64) error {
	return s.categoryRepo.UpdateUnitPrice(id, unitPrice)
}

// UpdateCategory 显式更新分类字段
func (s *CategoryService) UpdateCategory(id uint, updates map[string]interface{}) error {
	return s.categoryRepo.UpdateFields(id, updates)
}

// Delete 软删除分类
func (s *CategoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}
