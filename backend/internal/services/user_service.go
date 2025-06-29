package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

// UserService 用户服务 (不再使用接口)
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Create 创建用户
func (s *UserService) Create(user *models.User) error {
	// Hash password before saving
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return s.userRepo.Create(user)
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// GetAll 获取所有用户
func (s *UserService) GetAll() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// UpdatePassword 显式更新用户密码
func (s *UserService) UpdatePassword(id uint, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(id, hashedPassword)
}

// UpdateRealName 显式更新真实姓名
func (s *UserService) UpdateRealName(id uint, realName string) error {
	return s.userRepo.UpdateRealName(id, realName)
}

// UpdateRole 显式更新用户角色
func (s *UserService) UpdateRole(id uint, role string) error {
	return s.userRepo.UpdateRole(id, role)
}

// UpdateUser 显式更新用户字段
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	// 如果更新密码，需要先加密
	if password, exists := updates["password"]; exists {
		if passwordStr, ok := password.(string); ok {
			hashedPassword, err := HashPassword(passwordStr)
			if err != nil {
				return err
			}
			updates["password"] = hashedPassword
		}
	}
	return s.userRepo.UpdateFields(id, updates)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
