package services

import (
	"cms/models"

	"gorm.io/gorm"
)

type (
	UserService interface {
		GetUsers() ([]*models.User, error)
	}
	userService struct {
		db *gorm.DB
	}
)

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
