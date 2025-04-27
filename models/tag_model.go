package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:char(36)"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Description string    `json:"description"`

	Articles []*Article `json:"articles" gorm:"many2many:article_tags"`

	CommonModel
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
