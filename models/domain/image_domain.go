package domain

import "cms/models"

type (
	// 添加图片参数
	CreateImageParams struct {
		Title string `json:"title" validate:"required"`
		Hash  uint64 `json:"hash" validate:"required"`
	}

	// 添加图片响应
	CreateImageResponse []models.Image
)
