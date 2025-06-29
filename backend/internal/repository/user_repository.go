package repository

import (
	"battery-erp-backend/internal/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据仓库 (不再是接口实现)
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll 获取所有活跃用户
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// UpdatePassword 显式更新用户密码
func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// UpdateRealName 显式更新真实姓名
func (r *UserRepository) UpdateRealName(id uint, realName string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("real_name", realName).Error
}

// UpdateRole 显式更新用户角色
func (r *UserRepository) UpdateRole(id uint, role string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("role", role).Error
}

// UpdateFields 显式更新指定字段
func (r *UserRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除用户 (设置为非活跃状态)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}
