package models

type Image struct {
	ID    uint   `json:"id" gorm:"primarykey"`
	Title string `json:"title"`
	Hash  uint64 `json:"hash,string" gorm:"uniqueIndex;not null"`

	Articles []Article `json:"articles,omitempty" gorm:"many2many:article_images;"`

	CommonModel
}
