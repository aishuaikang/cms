package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:char(36)"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Alias       string    `json:"alias" gorm:"not null;unique"`
	Description string    `json:"description"`
	Sort        uint      `json:"sort" gorm:"not null;default:0"`

	Articles []Article `json:"articles"`

	CommonModel
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
