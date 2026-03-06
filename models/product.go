package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	IsExclusive bool   `json:"is_exclusive" gorm:"default:false;index"`

	CatalogID *uint    `json:"catalog_id" gorm:"index"`
	Catalog   *Catalog `json:"catalog,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
