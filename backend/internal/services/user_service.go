package services

import (
	"battery-erp-backend/internal/models"
	"battery-erp-backend/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(user *models.User) error {
	// Hash password before saving
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	
	return s.userRepo.Create(user)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}