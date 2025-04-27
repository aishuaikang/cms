package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleStatus uint8

const (
	StatusDraft ArticleStatus = iota
	StatusPublished
)

type Article struct {
	ID          uuid.UUID     `json:"id" gorm:"primary_key;type:char(36)"`
	Title       string        `json:"title" gorm:"not null;unique"`
	Description string        `json:"description"`
	Content     string        `json:"content"`
	Status      ArticleStatus `json:"status"`

	CategoryID uuid.UUID `json:"category_id"`

	Images []*Image `json:"images" gorm:"many2many:article_images"`

	Tags []*Tag `json:"tags" gorm:"many2many:article_tags"`

	UserID uuid.UUID `json:"user_id"`

	CommonModel
}

func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
