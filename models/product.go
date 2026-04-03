package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsExclusive bool           `json:"is_exclusive" gorm:"default:false;index"`
	CatalogID   *uint          `json:"catalog_id" gorm:"index"`
	Catalog     *Catalog       `json:"catalog,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Images      []ProductImage `json:"-" gorm:"foreignKey:ProductID"`
	ImageURL    string         `json:"imageUrl" gorm:"-"`
}
