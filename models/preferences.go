package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Preferences struct {
	gorm.Model
	Style            int16               `json:"style" gorm:"default:1"`
	FavoriteColors   pgtype.Array[int16] `json:"favoriteColors" gorm:"type:smallint[]"`
	PreferredBra     int16               `json:"preferredBra" gorm:"default:1"`
	PreferredModel   int16               `json:"preferredModel" gorm:"default:1"`
	PreferredPanties int16               `json:"preferredPanties" gorm:"default:1"`
	Size             int16               `json:"size" gorm:"default:1"`
	AllowedModels    string              `json:"allowedModels" gorm:"type:text"`
	NotAllowedModels string              `json:"notAllowedModels" gorm:"type:text"`
	Notes            string              `json:"notes" gorm:"type:text"`
}
