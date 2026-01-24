package models

import (
	"gorm.io/gorm"
)

type CatalogProduct struct {
	gorm.Model
	Price    float64 `json:"price"`
	IsBought bool    `json:"is_bought" gorm:"default:false"`

	CatalogID uint    `json:"catalog_id"`
	Catalog   Catalog `json:"catalog" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
