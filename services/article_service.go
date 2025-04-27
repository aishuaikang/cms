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
	// ErrArticleNotFound 文章不存在时
	ErrArticleNotFound = errors.New("文章不存在")
	// ErrArticleAlreadyExists 文章已存在时
	ErrArticleAlreadyExists = errors.New("文章已存在")
)

type (
	ArticleService interface {
		GetArticles() ([]*models.Article, error)
		CreateArticle(user_id uuid.UUID, article domain.CreateArticleParams) error
		UpdateArticle(id uuid.UUID, article domain.UpdateArticleParams) error
		DeleteArticle(id uuid.UUID) error
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
	if err := s.db.Preload(clause.Associations).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *articleService) CreateArticle(user_id uuid.UUID, params domain.CreateArticleParams) error {
	// 检查分类是否存在
	if err := s.db.Where("id = ?", params.CategoryID).First(&models.Category{}).Error; err != nil {
		return ErrCategoryNotFound
	}

	// 检查文章标题是否已存在
	if err := s.db.Where("title = ?", params.Title).First(&models.Article{}).Error; err == nil {
		return ErrArticleAlreadyExists
	}

	article := &models.Article{
		Title:       params.Title,
		Description: params.Description,
		Content:     params.Content,
		CategoryID:  params.CategoryID,
		Status:      params.Status,
		UserID:      user_id,
	}

	if err := s.db.Create(&article).Error; err != nil {
		return err
	}

	if params.ImageIds != nil {
		images := make([]*models.Image, 0)
		for _, id := range params.ImageIds {
			image := new(models.Image)
			if err := s.db.Where("id = ?", id).First(image).Error; err != nil {
				continue
			}
			images = append(images, image)
		}
		s.db.Model(&article).Association("Images").Append(images)
	}

	if params.TagIds != nil {
		tags := make([]*models.Tag, 0)
		for _, id := range params.TagIds {
			tag := new(models.Tag)
			// 检查标签是否存在
			if err := s.db.Where("id = ?", id).First(tag).Error; err != nil {
				continue // 如果标签不存在，跳过
			}
			tags = append(tags, tag)
		}
		s.db.Model(&article).Association("Tags").Append(tags)
	}

	return nil
}

func (s *articleService) UpdateArticle(id uuid.UUID, params domain.UpdateArticleParams) error {
	article := new(models.Article)
	// 检查分类是否存在
	if err := s.db.Where("id = ?", id).First(article).Error; err != nil {
		return ErrArticleNotFound
	}

	if params.Title != nil && article.Title != *params.Title {
		// 检查文章标题是否已存在
		if err := s.db.Where("title = ?", *params.Title).First(&models.Article{}).Error; err == nil {
			return ErrArticleAlreadyExists
		}

		article.Title = *params.Title
	}

	if params.Description != nil && article.Description != *params.Description {
		article.Description = *params.Description
	}

	if params.Content != nil && article.Content != *params.Content {
		article.Content = *params.Content
	}

	if params.CategoryID != nil && article.CategoryID != *params.CategoryID {
		// 检查分类是否存在
		if err := s.db.Where("id = ?", *params.CategoryID).First(&models.Category{}).Error; err != nil {
			return ErrCategoryNotFound
		}
		article.CategoryID = *params.CategoryID
	}

	if params.Status != nil && article.Status != *params.Status {
		article.Status = *params.Status
	}

	if params.ImageIds != nil {
		images := make([]*models.Image, 0)
		for _, id := range params.ImageIds {
			image := new(models.Image)
			// 检查图片是否存在
			if err := s.db.Where("id = ?", id).First(image).Error; err != nil {
				continue
			}
			images = append(images, image)
		}
		s.db.Model(&article).Association("Images").Replace(images)
	}

	if params.TagIds != nil {
		tags := make([]*models.Tag, 0)
		for _, id := range params.TagIds {
			tag := new(models.Tag)
			// 检查标签是否存在
			if err := s.db.Where("id = ?", id).First(tag).Error; err != nil {
				continue
			}
			tags = append(tags, tag)
		}
		s.db.Model(&article).Association("Tags").Replace(tags)
	}

	return s.db.Save(article).Error
}

func (s *articleService) DeleteArticle(id uuid.UUID) error {
	article := new(models.Article)
	// 检查文章是否存在
	if err := s.db.Where("id = ?", id).First(article).Error; err != nil {
		return ErrArticleNotFound
	}

	return s.db.Delete(article).Error
}
