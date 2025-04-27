package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key;type:char(36)"`
	Title string    `json:"title"`
	Hash  uint64    `json:"hash" gorm:"uniqueIndex;not null"`

	Articles []*Article `json:"articles" gorm:"many2many:article_images"`

	Users []*User `json:"users"`

	CommonModel
}

func (i *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}
