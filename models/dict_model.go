package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Dict struct {
	ID          uuid.UUID  `json:"id" gorm:"primary_key;type:char(36)"`
	Name        string     `json:"name" gorm:"not null;unique"`
	Code        string     `json:"code" gorm:"not null;unique"`
	Extra       string     `json:"extra" gorm:"type:mediumtext"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parentId"`
	ImageID     *uuid.UUID `json:"imageId"`

	CommonModel
}

func (d *Dict) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}
