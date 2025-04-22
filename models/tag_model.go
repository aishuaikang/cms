package models

type Tag struct {
	ID          uint   `json:"id,string" gorm:"primary_key"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`

	Articles []Article `json:"articles" gorm:"many2many:article_tags;"`

	CommonModel
}
