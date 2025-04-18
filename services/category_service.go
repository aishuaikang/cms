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
		CreateCategory(category domain.CreateCategoryParams) error
		UpdateCategory(id string, category domain.UpdateCategoryParams) error
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

func (s *categoryService) CreateCategory(category domain.CreateCategoryParams) error {
	// 是否存在
	if err := s.db.Where("name = ?", category.Name).First(&models.Category{}).Error; err == nil {
		return ErrCategoryAlreadyExists
	}

	categoryModel := &models.Category{
		Name:        category.Name,
		Description: category.Description,
	}

	return s.db.Create(&categoryModel).Error
}

func (s *categoryService) UpdateCategory(id string, category domain.UpdateCategoryParams) error {
	if err := s.db.Where("id = ?", id).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	categoryModel := &models.Category{}

	if category.Name != nil {
		categoryModel.Name = *category.Name
	}

	if category.Description != nil {
		categoryModel.Description = *category.Description
	}

	return s.db.Model(&models.Category{}).Where("id = ?", id).Updates(categoryModel).Error
}

func (s *categoryService) DeleteCategory(id string) error {
	if err := s.db.Where("id = ?", id).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	if err := s.db.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
