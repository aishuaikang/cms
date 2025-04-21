package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrCategoryNotFound 当分类不存在时
	ErrCategoryNotFound = errors.New("分类不存在")
	// ErrCategoryAlreadyExists 当分类已存在时
	ErrCategoryAlreadyExists = errors.New("分类已存在")
)

type (
	CategoryService interface {
		GetCategorys() ([]*models.Category, error)
		CreateCategory(params domain.CreateCategoryParams) error
		UpdateCategory(id string, params domain.UpdateCategoryParams) error
		DeleteCategory(id string) error
	}
	categoryService struct {
		db *gorm.DB
	}
)

func NewCategoryService(db *gorm.DB) CategoryService {
	return &categoryService{db: db}
}

func (s *categoryService) GetCategorys() ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) CreateCategory(params domain.CreateCategoryParams) error {
	// 检查分类名称是否已存在
	if err := s.db.Where("name = ?", params.Name).First(&models.Category{}).Error; err == nil {
		return ErrCategoryAlreadyExists
	}

	categoryModel := &models.Category{
		Name:        params.Name,
		Description: params.Description,
	}

	return s.db.Create(&categoryModel).Error
}

func (s *categoryService) UpdateCategory(id string, params domain.UpdateCategoryParams) error {
	// 检查分类是否存在
	if err := s.db.Where("id = ?", id).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	categoryModel := &models.Category{}

	if params.Name != nil {
		categoryModel.Name = *params.Name
	}

	if params.Description != nil {
		categoryModel.Description = *params.Description
	}

	return s.db.Model(&models.Category{}).Where("id = ?", id).Updates(categoryModel).Error
}

func (s *categoryService) DeleteCategory(id string) error {
	// 检查分类是否存在
	if err := s.db.Where("id = ?", id).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	if err := s.db.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
