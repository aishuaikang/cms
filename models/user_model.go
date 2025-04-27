package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key;type:char(36)"`
	Nickname string    `json:"nickname" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"uniqueIndex;not null"`
	Username string    `json:"username" gorm:"uniqueIndex;not null"`
	Password string    `json:"-" gorm:"not null"`
	IsSuper  bool      `json:"is_super" gorm:"not null"`

	ImageID *uuid.UUID `json:"image_id"`

	Articles []*Article `json:"articles"`

	CommonModel
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
