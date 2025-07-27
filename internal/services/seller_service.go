package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

// SellerService 卖家服务
type SellerService struct {
	repo *repository.SellerRepository
}

// NewSellerService 创建卖家服务实例
func NewSellerService(repo *repository.SellerRepository) *SellerService {
	return &SellerService{repo: repo}
}

// Create 创建卖家
func (s *SellerService) Create(seller *models.Seller) error {
	return s.repo.Create(seller)
}

// Update 更新卖家信息
func (s *SellerService) Update(id uint, updates map[string]interface{}) error {
	return s.repo.Update(id, updates)
}

// Delete 删除卖家
func (s *SellerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetByID 根据ID获取卖家
func (s *SellerService) GetByID(id uint) (*models.Seller, error) {
	return s.repo.GetByID(id)
}

// GetAll 获取所有卖家
func (s *SellerService) GetAll() ([]models.Seller, error) {
	return s.repo.GetAll()
}

// GetByName 根据名称获取卖家
func (s *SellerService) GetByName(name string) ([]models.Seller, error) {
	return s.repo.GetByName(name)
}
