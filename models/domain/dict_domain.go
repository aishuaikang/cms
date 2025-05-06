package domain

import "github.com/google/uuid"

type (
	// 添加字典参数
	CreateDictParams struct {
		Name        string     `json:"name" validate:"required"`
		Code        string     `json:"code" validate:"required"`
		Extra       string     `json:"extra"`
		Description string     `json:"description"`
		ParentID    *uuid.UUID `json:"parentId"`
		ImageID     *uuid.UUID `json:"imageId"`
	}
	// 修改字典参数
	UpdateDictParams struct {
		Name        *string    `json:"name"`
		Code        *string    `json:"code"`
		Extra       *string    `json:"extra"`
		Description *string    `json:"description"`
		ImageID     *uuid.UUID `json:"imageId"`

		// ParentID    *uint   `json:"parent_id"`
	}
)
