package models

type Image struct {
	ID    uint   `json:"id,string" gorm:"primarykey"`
	Title string `json:"title"`
	Hash  uint64 `json:"hash,string" gorm:"uniqueIndex;not null"`

	Articles []Article `json:"articles" gorm:"many2many:article_images;"`

	Users []User `json:"users"`

	CommonModel
}
