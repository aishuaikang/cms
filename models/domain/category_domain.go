package domain

type (
	// 添加分类参数
	CreateCategoryParams struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}
	// 修改分类参数
	UpdateCategoryParams struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
)
