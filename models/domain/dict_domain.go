package domain

type (
	// 添加字典参数
	CreateDictParams struct {
		Name        string `json:"name" validate:"required"`
		Code        string `json:"code" validate:"required"`
		Extra       string `json:"extra"`
		Description string `json:"description"`
		ParentID    *uint  `json:"parent_id,string"`
	}
	// 修改字典参数
	UpdateDictParams struct {
		Name        *string `json:"name"`
		Code        *string `json:"code"`
		Extra       *string `json:"extra"`
		Description *string `json:"description"`
		ParentID    *uint   `json:"parent_id,string"`
	}
)
