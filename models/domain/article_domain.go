package domain

import (
	"cms/models"

	"github.com/google/uuid"
)

type (
	// 获取文章列表参数
	GetArticleListParams struct {
		Page       int                   `json:"page" validate:"required,min=1"`
		PageSize   int                   `json:"pageSize" validate:"required,min=1,max=100"`
		Title      *string               `json:"title"`
		Status     *models.ArticleStatus `json:"status" validate:"omitempty,oneof=0 1"`
		CategoryID *uuid.UUID            `json:"categoryId"`
	}

	// 添加文章参数
	CreateArticleParams struct {
		Title       string               `json:"title" validate:"required"`
		Description string               `json:"description" validate:"required"`
		Content     string               `json:"content" validate:"required"`
		CategoryID  uuid.UUID            `json:"categoryId" validate:"required"`
		Status      models.ArticleStatus `json:"status" validate:"oneof=0 1"`
		ImageIds    []uuid.UUID          `json:"imageIds"`
		TagIds      []uuid.UUID          `json:"tagIds"`
	}
	// 修改文章参数
	UpdateArticleParams struct {
		Title       *string               `json:"title"`
		Description *string               `json:"description"`
		Content     *string               `json:"content"`
		CategoryID  *uuid.UUID            `json:"categoryId"`
		Status      *models.ArticleStatus `json:"status"`
		ImageIds    []uuid.UUID           `json:"imageIds"`
		TagIds      []uuid.UUID           `json:"tagIds"`
	}

	// 获取文章列表返回值
	GetArticlesByCategoryAliasWithCacheParams struct {
		Page     int `json:"page" validate:"required,min=1"`
		PageSize int `json:"pageSize" validate:"required,min=1,max=100"`
	}

	// GetRelatedArticlesByIDWithCacheParams 获取相关文章参数
	GetRelatedArticlesByIDWithCacheParams struct {
		Page     int `json:"page" validate:"required,min=1"`
		PageSize int `json:"pageSize" validate:"required,min=1,max=100"`
	}
)
