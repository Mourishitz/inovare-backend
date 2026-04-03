package models

import "gorm.io/gorm"

type ProductImage struct {
	gorm.Model
	ProductID uint   `json:"product_id" gorm:"index;not null"`
	ImageURL  string `json:"image_url" gorm:"not null"`
	IsPrimary bool   `json:"is_primary" gorm:"default:false"`
}
