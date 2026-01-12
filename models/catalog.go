package models

import (
	"gorm.io/gorm"
)

type Catalog struct {
	gorm.Model
	URL      string `json:"url" gorm:"type:text"`
	Approved bool   `json:"approved" gorm:"default:false"`
}
