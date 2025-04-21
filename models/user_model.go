package models

type User struct {
	ID       uint   `json:"id,string" gorm:"primarykey"`
	Nickname string `json:"nickname" gorm:"not null"`
	Phone    string `json:"phone" gorm:"uniqueIndex;not null"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	IsSuper  bool   `json:"is_super" gorm:"not null"`

	ImageID *uint `json:"image_id,string"`

	CommonModel
}
