package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrCategoryNotFound 分类不存在时
	ErrCategoryNotFound = errors.New("分类不存在")
	// ErrCategoryNameAlreadyExists 分类名称已存在
	ErrCategoryNameAlreadyExists = errors.New("分类名称已存在")
	// ErrCategoryAliasAlreadyExists 分类别名已存在
	ErrCategoryAliasAlreadyExists = errors.New("分类别名已存在")
	// ErrCategoryHasArticles 分类下存在文章
	ErrCategoryHasArticles = errors.New("分类下存在文章")
)

type (
	CategoryService interface {
		GetCategorys() ([]*models.Category, error)
		CreateCategory(params domain.CreateCategoryParams) error
		UpdateCategory(id uuid.UUID, params domain.UpdateCategoryParams) error
		DeleteCategory(id uuid.UUID) error
		GetCategorysWithCache() ([]*models.Category, error)
		GetCategoryByAliasWithCache(alias string) (*models.Category, error)
		GetCategoryByIDWithCache(articleID uuid.UUID) (*models.Category, error)
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
	if err := s.db.Order("sort DESC, created_at ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) CreateCategory(params domain.CreateCategoryParams) error {
	// 检查分类名称是否已存在
	if err := s.db.Where("name = ?", params.Name).First(&models.Category{}).Error; err == nil {
		return ErrCategoryNameAlreadyExists
	}

	// 检查分类别名是否已存在
	if err := s.db.Where("alias = ?", params.Alias).First(&models.Category{}).Error; err == nil {
		return ErrCategoryAliasAlreadyExists
	}

	categoryModel := &models.Category{
		Name:        params.Name,
		Alias:       params.Alias,
		Sort:        params.Sort,
		Description: params.Description,
	}

	return s.db.Create(&categoryModel).Error
}

func (s *categoryService) UpdateCategory(id uuid.UUID, params domain.UpdateCategoryParams) error {
	category := new(models.Category)

	// 检查分类是否存在
	if err := s.db.Where("id = ?", id).First(category).Error; err != nil {
		return ErrCategoryNotFound
	}

	if params.Name != nil && category.Name != *params.Name {
		// 检查分类名称是否已存在
		if err := s.db.Where("name = ?", *params.Name).First(&models.Category{}).Error; err == nil {
			return ErrCategoryNameAlreadyExists
		}
		category.Name = *params.Name
	}

	if params.Alias != nil && category.Alias != *params.Alias {
		// 检查分类别名是否已存在
		if err := s.db.Where("alias = ?", *params.Alias).First(&models.Category{}).Error; err == nil {
			return ErrCategoryAliasAlreadyExists
		}
		category.Alias = *params.Alias
	}

	if params.Sort != nil && category.Sort != *params.Sort {
		category.Sort = *params.Sort
	}

	if params.Description != nil && category.Description != *params.Description {
		category.Description = *params.Description
	}

	return s.db.Save(category).Error
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	category := new(models.Category)
	// 检查分类是否存在
	if err := s.db.Preload(clause.Associations).Where("id = ?", id).First(category).Error; err != nil {
		return ErrCategoryNotFound
	}

	if len(category.Articles) > 0 {
		return ErrCategoryHasArticles
	}

	return s.db.Delete(category).Error
}

func (s *categoryService) GetCategorysWithCache() ([]*models.Category, error) {
	var categories []*models.Category
	if err := s.db.Order("sort DESC, created_at ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) GetCategoryByAliasWithCache(alias string) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("alias = ?", alias).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categoryService) GetCategoryByIDWithCache(articleID uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := s.db.Where("id = ?", articleID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
