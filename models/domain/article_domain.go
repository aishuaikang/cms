package domain

import "cms/models"

type (
	// 添加文章参数
	CreateArticleParams struct {
		Title       string               `json:"title" validate:"required"`
		Description string               `json:"description" validate:"required"`
		Content     string               `json:"content" validate:"required"`
		CategoryID  uint                 `json:"category_id,string" validate:"required"`
		Status      models.ArticleStatus `json:"status" validate:"required,oneof=0 1"`
		ImageIds    []uint               `json:"image_ids"`
		TagIds      []uint               `json:"tag_ids"`
	}
	// 修改文章参数
	UpdateArticleParams struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Content     *string `json:"content"`
		CategoryID  *uint   `json:"category_id,string"`
		ImageIds    []uint  `json:"image_ids"`
		TagIds      []uint  `json:"tag_ids"`
	}
)
