package models

import (
	"encoding/json"
	"inovare-backend/config"
	"inovare-backend/models/enums"

	"gorm.io/gorm"
)

type Catalog struct {
	gorm.Model
	URL      string `json:"-" gorm:"type:text;unique;not null"` // Store unique ID only
	Package  int16  `json:"-" gorm:"not null"`
	Approved bool   `json:"approved" gorm:"default:false"`
}

// MarshalJSON implements custom JSON marshaling for Catalog
func (c Catalog) MarshalJSON() ([]byte, error) {
	type Alias Catalog
	return json.Marshal(&struct {
		*Alias
		URL     string `json:"url"`
		Package string `json:"package"`
	}{
		Alias:   (*Alias)(&c),
		URL:     config.GetConfig().FrontendURL + "/" + c.URL,
		Package: getPackageName(c.Package),
	})
}

func getPackageName(id int16) string {
	if name, ok := enums.PackageNames[id]; ok {
		return name
	}
	return "Unknown"
}
