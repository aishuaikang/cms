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
		UpdateUser(id uint, params domain.UpdateUserParams) error
		DeleteUser(id uint) error
		Login(params domain.LoginParams) (*models.User, error)
		CreateInitialUser(initialUsername, initialPassword string) error
		// 根据用户id获取是否是超级管理员
		GetUserIsSuper(id uint) (bool, error)
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

func (s *userService) UpdateUser(id uint, params domain.UpdateUserParams) error {
	user := new(models.User)
	// 检查用户是否存在
	if err := s.db.Where("id = ?", id).First(user).Error; err != nil {
		return ErrUserNotFound
	}

	if params.Username != nil && user.Username != *params.Username {
		// 检查用户名是否已存在
		if err := s.db.Where("username = ?", *params.Username).First(&models.User{}).Error; err == nil {
			return ErrUsernameExists
		}

		user.Username = *params.Username
	}

	if params.Phone != nil && user.Phone != *params.Phone {
		// 检查手机号是否已存在
		if err := s.db.Where("phone = ?", *params.Phone).First(&models.User{}).Error; err == nil {
			return ErrPhoneExists
		}
		user.Phone = *params.Phone
	}

	if params.Nickname != nil && user.Nickname != *params.Nickname {
		user.Nickname = *params.Nickname
	}

	if params.Password != nil {
		// 对密码进行加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	if params.Nickname != nil && user.Nickname != *params.Nickname {
		user.Nickname = *params.Nickname
	}

	if params.ImageID != nil && user.ImageID != params.ImageID {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", *params.ImageID).First(&models.Image{}).Error; err != nil {
			return ErrImageNotFound
		}
		user.ImageID = params.ImageID
	}

	return s.db.Save(user).Error
}

func (s *userService) DeleteUser(id uint) error {
	user := new(models.User)
	// 检查用户是否存在
	if err := s.db.Where("id = ?", id).First(user).Error; err != nil {
		return ErrUserNotFound
	}

	return s.db.Delete(user).Error
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
		Nickname: "ask",
		Phone:    "18333435634",
		Username: initialUsername,
		Password: string(hashedPassword),
		IsSuper:  true,
	}

	return s.db.Create(user).Error
}

func (s *userService) GetUserIsSuper(id uint) (bool, error) {
	user := new(models.User)
	// 检查用户是否存在
	if err := s.db.Where("id = ?", id).First(user).Error; err != nil {
		return false, ErrUserNotFound
	}

	return user.IsSuper, nil
}
