package scopes

import (
	"cms/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ArticleScope interface {
		Title(title *string) func(*gorm.DB) *gorm.DB
		Category(categoryID *uuid.UUID) func(*gorm.DB) *gorm.DB
		Status(status *models.ArticleStatus) func(*gorm.DB) *gorm.DB
	}
	articleScope struct {
		db *gorm.DB
	}
)

func NewArticleScope(db *gorm.DB) ArticleScope {
	return &articleScope{db: db}
}

func (s *articleScope) Title(title *string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if title != nil {
			return db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", *title))
		}
		return db
	}
}

func (s *articleScope) Category(categoryID *uuid.UUID) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if categoryID != nil {
			return db.Where("category_id = ?", *categoryID)
		}
		return db
	}
}

func (s *articleScope) Status(status *models.ArticleStatus) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != nil {
			return db.Where("status = ?", *status)
		}
		return db
	}
}
