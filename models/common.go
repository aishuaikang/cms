package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CommonModel struct {
	CreatedAt CustomTime     `json:"createdAt"`
	UpdatedAt CustomTime     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CommonNotDeletedModel struct {
	CreatedAt CustomTime `json:"createdAt"`
	UpdatedAt CustomTime `json:"updatedAt"`
}

type CustomTime time.Time

func (ct *CustomTime) Scan(value interface{}) error {
	if t, ok := value.(time.Time); ok {
		*ct = CustomTime(t)
		return nil
	}
	return fmt.Errorf("failed to scan CustomTime: %v", value)
}

func (ct CustomTime) Value() (driver.Value, error) {
	t := time.Time(ct)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	if t.IsZero() {
		return json.Marshal(nil)
	}
	return json.Marshal(t.Format("2006-01-02 15:04:05"))
}
