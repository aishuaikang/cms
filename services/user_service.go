package services

import (
	"cms/models"
	"cms/models/domain"
	"cms/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	// ErrUsernameExists 用户名已存在
	ErrUsernameExists = errors.New("用户名已存在")
	// ErrPhoneExists 手机号已存在
	ErrPhoneExists = errors.New("手机号已存在")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("用户不存在")
	// ErrUsernameNotFound 用户名不存在
	ErrUsernameNotFound = errors.New("用户名不存在")
	// ErrPasswordIncorrect 密码错误
	ErrPasswordIncorrect = errors.New("密码错误")
)

type (
	UserService interface {
		GetUsers() ([]*models.User, error)
		CreateUser(params domain.CreateUserParams) error
		UpdateUser(id string, params domain.UpdateUserParams) error
		DeleteUser(id string) error
		Login(params domain.LoginParams) (*models.User, error)
		CreateInitialUser(initialUsername, initialPassword string) error
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

func (s *userService) CreateUser(user domain.CreateUserParams) error {
	// 检查用户名是否存在
	if err := s.db.Where("username = ?", user.Username).First(&models.User{}).Error; err == nil {
		return ErrUsernameExists
	}

	// 检查手机号是否存在
	if err := s.db.Where("phone = ?", user.Phone).First(&models.User{}).Error; err == nil {
		return ErrPhoneExists
	}

	// 对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	categoryModel := &models.User{
		Nickname: user.Nickname,
		Phone:    user.Phone,
		Username: user.Username,
		Password: string(hashedPassword),
	}

	if user.ImageID != nil {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", *user.ImageID).First(&models.Image{}).Error; err != nil {
			return ErrImageNotFound
		}
		categoryModel.ImageID = user.ImageID
	}

	return s.db.Create(&categoryModel).Error
}

func (s *userService) UpdateUser(id string, params domain.UpdateUserParams) error {
	// 检查用户是否存在
	if err := s.db.Where("id = ?", id).First(&models.User{}).Error; err != nil {
		return ErrUserNotFound
	}

	userModel := &models.User{}

	// 检查用户名是否已存在
	if params.Username != nil {
		if err := s.db.Where("username = ?", *params.Username).First(&models.User{}).Error; err == nil {
			return ErrUsernameExists
		}

		userModel.Username = *params.Username
	}

	// 检查手机号是否已存在
	if params.Phone != nil {
		if err := s.db.Where("phone = ?", *params.Phone).First(&models.User{}).Error; err == nil {
			return ErrPhoneExists
		}
		userModel.Phone = *params.Phone
	}

	if params.Nickname != nil {
		userModel.Nickname = *params.Nickname
	}

	if params.Password != nil {
		// 对密码进行加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userModel.Password = string(hashedPassword)
	}

	if params.Nickname != nil {
		userModel.Nickname = *params.Nickname
	}

	if params.ImageID != nil {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", *params.ImageID).First(&models.Image{}).Error; err != nil {
			return ErrImageNotFound
		}
		userModel.ImageID = params.ImageID
	}

	return s.db.Model(&models.User{}).Where("id = ?", id).Updates(userModel).Error
}

func (s *userService) DeleteUser(id string) error {
	// 检查用户是否存在
	if err := s.db.Where("id = ?", id).First(&models.User{}).Error; err != nil {
		return ErrUserNotFound
	}

	if err := s.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(params domain.LoginParams) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", params.Username).First(&user).Error; err != nil {
		return nil, ErrUsernameNotFound
	}

	// 验证密码
	if err := utils.VerifyPassword(user.Password, params.Password); err != nil {
		return nil, ErrPasswordIncorrect
	}

	return &user, nil
}

func (s *userService) CreateInitialUser(initialUsername, initialPassword string) error {
	// 检查是否已经存在用户
	var count int64
	if err := s.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 创建初始用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(initialPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Nickname: "admin",
		Phone:    "12345678901",
		Username: initialUsername,
		Password: string(hashedPassword),
		IsSuper:  true,
	}

	return s.db.Create(user).Error
}
