package models

type Category struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`

	// Articles []Article `json:"articles"`

	CommonModel
}
