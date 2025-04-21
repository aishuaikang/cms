package domain

type (
	// 添加文章参数
	CreateArticleParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Content     string `json:"content" validate:"required"`
		CategoryID  uint   `json:"category_id,string" validate:"required"`
		ImageIds    []uint `json:"image_ids"`
	}
	// 修改文章参数
	UpdateArticleParams struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Content     *string `json:"content"`
		CategoryID  *uint   `json:"category_id,string"`
	}
)
