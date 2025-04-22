package models

type Article struct {
	ID          uint   `json:"id,string" gorm:"primary_key"`
	Title       string `json:"title" gorm:"not null;unique"`
	Description string `json:"description"`
	Content     string `json:"content"`

	CategoryID uint `json:"category_id,string"`

	Images []Image `json:"images" gorm:"many2many:article_images;"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tags;"`

	UserID uint `json:"user_id,string"`

	CommonModel
}
