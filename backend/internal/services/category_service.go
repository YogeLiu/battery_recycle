package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

type categoryService struct {
	categoryRepo  repository.CategoryRepository
	inventoryRepo repository.InventoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository, inventoryRepo repository.InventoryRepository) CategoryService {
	return &categoryService{
		categoryRepo:  categoryRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *categoryService) Create(category *models.BatteryCategory) error {
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

func (s *categoryService) GetByID(id uint) (*models.BatteryCategory, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetAll() ([]models.BatteryCategory, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) Update(category *models.BatteryCategory) error {
	return s.categoryRepo.Update(category)
}

func (s *categoryService) Delete(id uint) error {
	// 软删除分类
	return s.categoryRepo.Delete(id)
}
