package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
	"time"

	"gorm.io/gorm"
)

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	categoryRepo  repository.CategoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository, categoryRepo repository.CategoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		categoryRepo:  categoryRepo,
	}
}

func (s *inventoryService) GetByCategoryID(categoryID uint) (*models.Inventory, error) {
	// 先检查分类是否存在
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventoryRepo.GetByCategoryID(categoryID)
	if err != nil {
		// 使用正确的gorm错误判断
		if err == gorm.ErrRecordNotFound {
			// 如果库存记录不存在，创建一个新的
			newInventory := &models.Inventory{
				CategoryID:      categoryID,
				CurrentWeightKg: 0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			if createErr := s.inventoryRepo.Create(newInventory); createErr != nil {
				return nil, createErr
			}
			return s.inventoryRepo.GetByCategoryID(categoryID)
		}
		return nil, err
	}
	return inventory, nil
}

func (s *inventoryService) GetAll() ([]models.Inventory, error) {
	return s.inventoryRepo.GetAll()
}

// InitializeInventoryForCategory 为新分类初始化库存记录
func (s *inventoryService) InitializeInventoryForCategory(categoryID uint) error {
	// 检查是否已存在库存记录
	_, err := s.inventoryRepo.GetByCategoryID(categoryID)
	if err == nil {
		return nil // 已存在，无需初始化
	}

	if err != gorm.ErrRecordNotFound {
		return err // 其他错误
	}

	// 创建新的库存记录
	inventory := &models.Inventory{
		CategoryID:      categoryID,
		CurrentWeightKg: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return s.inventoryRepo.Create(inventory)
}
