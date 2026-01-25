package requests

import "time"

type CreateShowerRequest struct {
	Guests      uint      `json:"guests" binding:"required,min=1"`
	ShowerDate  time.Time `json:"shower_date" binding:"required"`
	WeddingDate time.Time `json:"wedding_date" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	HostID      uint      `json:"-"`
}

type UpdateShowerRequest struct {
	Guests      *uint      `json:"guests" binding:"omitempty,min=1"`
	ShowerDate  *time.Time `json:"shower_date" binding:"omitempty"`
	WeddingDate *time.Time `json:"wedding_date" binding:"omitempty"`
	Location    *string    `json:"location" binding:"omitempty"`
}

type AddCatalogRequest struct {
	Package int16 `json:"package" binding:"required,oneof=1 2"`
}

type AddPreferencesRequest struct {
	Style            int16   `json:"style" binding:"required"`
	FavoriteColors   []int16 `json:"favoriteColors" binding:"required"`
	PreferredBra     int16   `json:"preferredBra" binding:"required"`
	PreferredModel   int16   `json:"preferredModel" binding:"required"`
	PreferredPanties int16   `json:"preferredPanties" binding:"required"`
	Size             int16   `json:"size" binding:"required"`
	AllowedModels    string  `json:"allowedModels" binding:"required"`
	NotAllowedModels string  `json:"notAllowedModels" binding:"required"`
	Notes            string  `json:"notes" binding:"required"`
}
