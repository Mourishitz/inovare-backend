package models

import (
	"time"

	"gorm.io/gorm"
)

type Shower struct {
	gorm.Model
	Guests      uint      `json:"guests"`
	ShowerDate  time.Time `json:"shower_date"`
	WeddingDate time.Time `json:"wedding_date"`
	Location    string    `json:"location"`

	HostID uint `json:"host_id"`
	Host   User `json:"host" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	CatalogID *uint    `json:"catalog_id"`
	Catalog   *Catalog `json:"catalog,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	PreferencesID *uint        `json:"preferences_id"`
	Preferences   *Preferences `json:"preferences,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
