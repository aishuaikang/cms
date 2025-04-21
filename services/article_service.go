package services

import (
	"cms/models"
	"cms/models/domain"
	"errors"

	"gorm.io/gorm"
)

var (
	// ErrArticleNotFound 当文章不存在时
	ErrArticleNotFound = errors.New("文章不存在")
	// ErrArticleAlreadyExists 当文章已存在时
	ErrArticleAlreadyExists = errors.New("文章已存在")
)

type (
	ArticleService interface {
		GetArticles() ([]*models.Article, error)
		CreateArticle(article domain.CreateArticleParams) error
		UpdateArticle(id string, article domain.UpdateArticleParams) error
		DeleteArticle(id string) error
	}
	articleService struct {
		db *gorm.DB
	}
)

func NewArticleService(db *gorm.DB) ArticleService {
	return &articleService{db: db}
}

func (s *articleService) GetArticles() ([]*models.Article, error) {
	var categories []*models.Article
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *articleService) CreateArticle(article domain.CreateArticleParams) error {
	// 检查分类是否存在
	if err := s.db.Where("id = ?", article.CategoryID).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	// 检查文章标题是否已存在
	if err := s.db.Where("title = ?", article.Title).First(&models.Article{}).Error; err == nil {
		return ErrArticleAlreadyExists
	}

	images := make([]models.Image, 0)

	for _, id := range article.ImageIds {
		// 检查图片是否存在
		if err := s.db.Where("id = ?", id).First(&models.Image{}).Error; err != nil {
			continue
		}
		images = append(images, models.Image{
			ID: id,
		})
	}

	articleModel := &models.Article{
		Title:       article.Title,
		Description: article.Description,
		Content:     article.Content,
		CategoryID:  article.CategoryID,
		Images:      images,
	}

	return s.db.Create(&articleModel).Error
}

func (s *articleService) UpdateArticle(id string, article domain.UpdateArticleParams) error {
	// 检查分类是否存在
	if err := s.db.Where("id = ?", id).First(&models.Article{}).Error; err != nil {
		return ErrArticleNotFound
	}

	articleModel := &models.Article{}

	if article.Title != nil {
		// 检查文章标题是否已存在
		if err := s.db.Where("title = ?", *article.Title).First(&models.Article{}).Error; err == nil {
			return ErrArticleAlreadyExists
		}

		articleModel.Title = *article.Title
	}

	if article.Description != nil {
		articleModel.Description = *article.Description
	}

	if article.Content != nil {
		articleModel.Content = *article.Content
	}

	if article.CategoryID != nil {
		// 检查分类是否存在
		if err := s.db.Where("id = ?", *article.CategoryID).First(&models.Category{}).Error; err != nil {
			return ErrCategoryNotFound
		}
		articleModel.CategoryID = *article.CategoryID
	}

	if article.Content != nil {
		articleModel.Content = *article.Content
	}

	return s.db.Model(&models.Article{}).Where("id = ?", id).Updates(articleModel).Error
}

func (s *articleService) DeleteArticle(id string) error {
	// 检查文章是否存在
	if err := s.db.Where("id = ?", id).First(&models.Article{}).Error; err != nil {
		return ErrArticleNotFound
	}

	if err := s.db.Delete(&models.Article{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
