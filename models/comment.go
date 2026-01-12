package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `json:"content" gorm:"type:text"`

	AuthorID uint `json:"author_id"`
	Author   User `json:"author" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CatalogID uint    `json:"catalog_id"`
	Catalog   Catalog `json:"catalog" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
