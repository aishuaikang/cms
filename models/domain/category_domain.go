package domain

type (
	// 添加分类参数
	CreateCategoryParams struct {
		Name  string `json:"name" validate:"required"`
		Alias string `json:"alias" validate:"required"`
		Sort  uint   `json:"sort" validate:"required"`

		Description string `json:"description"`
	}
	// 修改分类参数
	UpdateCategoryParams struct {
		Name        *string `json:"name"`
		Alias       *string `json:"alias"`
		Sort        *uint   `json:"sort"`
		Description *string `json:"description"`
	}
)
