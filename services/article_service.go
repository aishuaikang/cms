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
		CreateArticle(user_id uint, article domain.CreateArticleParams) error
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

func (s *articleService) CreateArticle(user_id uint, article domain.CreateArticleParams) error {

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
		image := &models.Image{}
		// 检查图片是否存在
		if err := s.db.Where("id = ?", id).First(image).Error; err != nil {
			continue
		}
		images = append(images, *image)
	}

	tags := make([]models.Tag, 0)
	for _, id := range article.TagIds {
		tag := &models.Tag{}
		// 检查标签是否存在
		if err := s.db.Where("id = ?", id).First(tag).Error; err != nil {
			continue
		}
		tags = append(tags, *tag)
	}

	articleModel := &models.Article{
		Title:       article.Title,
		Description: article.Description,
		Content:     article.Content,
		CategoryID:  article.CategoryID,
		Images:      images,
		Tags:        tags,
		UserID:      user_id,
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

	if article.ImageIds != nil {
		images := make([]models.Image, 0)
		for _, id := range article.ImageIds {
			image := &models.Image{}
			// 检查图片是否存在
			if err := s.db.Where("id = ?", id).First(image).Error; err != nil {
				continue
			}
			images = append(images, *image)
		}
		articleModel.Images = images
	}

	if article.TagIds != nil {
		tags := make([]models.Tag, 0)
		for _, id := range article.TagIds {
			tag := &models.Tag{}
			// 检查标签是否存在
			if err := s.db.Where("id = ?", id).First(tag).Error; err != nil {
				continue
			}
			tags = append(tags, *tag)
		}
		articleModel.Tags = tags
	}

	return s.db.Model(&models.Article{}).Where("id = ?", id).Updates(articleModel).Error
}

func (s *articleService) DeleteArticle(id string) error {
	article := new(models.Article)
	// 检查文章是否存在
	if err := s.db.Where("id = ?", id).First(article).Error; err != nil {
		return ErrArticleNotFound
	}

	return s.db.Delete(article).Error
}
