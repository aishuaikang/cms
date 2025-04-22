package domain

type (
	// 添加标签参数
	CreateTagParams struct {
		Name        string  `json:"name" validate:"required"`
		Description *string `json:"description"`
	}
	// 修改标签参数
	UpdateTagParams struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
)
