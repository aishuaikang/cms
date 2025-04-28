package services

import (
	"cms/models"
	"cms/models/domain"
	"cms/models/scopes"
	"errors"
	"math"

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
		GetArticles(params domain.GetArticleListParams) (*domain.LimitResponse[*models.Article], error)
		CreateArticle(user_id uuid.UUID, article domain.CreateArticleParams) error
		UpdateArticle(id uuid.UUID, article domain.UpdateArticleParams) error
		DeleteArticle(id uuid.UUID) error
	}
	articleService struct {
		db           *gorm.DB
		articleScope scopes.ArticleScope
	}
)

func NewArticleService(db *gorm.DB, articleScope scopes.ArticleScope) ArticleService {
	return &articleService{db: db, articleScope: articleScope}
}

func (s *articleService) GetArticles(params domain.GetArticleListParams) (*domain.LimitResponse[*models.Article], error) {
	var count int64
	var articles []*models.Article

	// 基础查询
	model := s.db.Model(&models.Article{})
	// 统计总数
	if err := model.Scopes(
		s.articleScope.Title(params.Title),
		s.articleScope.Category(params.CategoryID),
		s.articleScope.Status(params.Status),
	).Count(&count).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := model.Scopes(
		s.articleScope.Title(params.Title),
		s.articleScope.Category(params.CategoryID),
		s.articleScope.Status(params.Status),
		scopes.PaginationScope(params.Page, params.PageSize),
	).Preload(clause.Associations).Find(&articles).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(count) / float64(params.PageSize)))

	return &domain.LimitResponse[*models.Article]{
		Total: count,
		Rows:  articles,
		Pages: totalPages,
	}, nil
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
