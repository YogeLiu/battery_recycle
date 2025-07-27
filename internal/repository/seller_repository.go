package repository

import (
	"battery-erp-backend/internal/models"

	"gorm.io/gorm"
)

// SellerRepository 卖家数据仓库
type SellerRepository struct {
	db *gorm.DB
}

// NewSellerRepository 创建卖家仓库实例
func NewSellerRepository(db *gorm.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

// Create 创建卖家
func (r *SellerRepository) Create(seller *models.Seller) error {
	return r.db.Create(seller).Error
}

// Update 更新卖家信息
func (r *SellerRepository) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.Seller{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除卖家 (软删除)
func (r *SellerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Seller{}, id).Error
}

// GetByID 根据ID获取卖家
func (r *SellerRepository) GetByID(id uint) (*models.Seller, error) {
	var seller models.Seller
	err := r.db.First(&seller, id).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// GetAll 获取所有卖家
func (r *SellerRepository) GetAll() ([]models.Seller, error) {
	var sellers []models.Seller
	err := r.db.Find(&sellers).Error
	return sellers, err
}

// GetByName 根据名称获取卖家
func (r *SellerRepository) GetByName(name string) ([]models.Seller, error) {
	var sellers []models.Seller
	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&sellers).Error
	return sellers, err
}
