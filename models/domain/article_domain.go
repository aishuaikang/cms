package domain

import (
	"cms/models"

	"github.com/google/uuid"
)

type (
	// 添加文章参数
	CreateArticleParams struct {
		Title       string               `json:"title" validate:"required"`
		Description string               `json:"description" validate:"required"`
		Content     string               `json:"content" validate:"required"`
		CategoryID  uuid.UUID            `json:"category_id" validate:"required"`
		Status      models.ArticleStatus `json:"status" validate:"required,oneof=0 1"`
		ImageIds    []uuid.UUID          `json:"image_ids"`
		TagIds      []uuid.UUID          `json:"tag_ids"`
	}
	// 修改文章参数
	UpdateArticleParams struct {
		Title       *string               `json:"title"`
		Description *string               `json:"description"`
		Content     *string               `json:"content"`
		CategoryID  *uuid.UUID            `json:"category_id"`
		Status      *models.ArticleStatus `json:"status"`
		ImageIds    []uuid.UUID           `json:"image_ids"`
		TagIds      []uuid.UUID           `json:"tag_ids"`
	}
)
