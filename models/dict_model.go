package models

type Dict struct {
	ID          uint   `json:"id,string" gorm:"primarykey"`
	Name        string `json:"name" gorm:"not null;unique"`
	Code        string `json:"code" gorm:"not null;unique"`
	Extra       string `json:"extra"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id,string"`

	CommonModel
}
