package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"inovare-backend/models/enums"
	"inovare-backend/models/enums/preferred"

	"gorm.io/gorm"
)

// Int16Array is a custom type for PostgreSQL int16 arrays
type Int16Array []int16

func (a *Int16Array) Scan(value interface{}) error {
	if value == nil {
		*a = []int16{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// PostgreSQL returns arrays as "{1,2,3}" format
		str := string(v)
		str = strings.Trim(str, "{}")
		if str == "" {
			*a = []int16{}
			return nil
		}

		parts := strings.Split(str, ",")
		result := make([]int16, len(parts))
		for i, part := range parts {
			var val int16
			_, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &val)
			if err != nil {
				return err
			}
			result[i] = val
		}
		*a = result
		return nil
	case string:
		// Handle string format
		str := strings.Trim(v, "{}")
		if str == "" {
			*a = []int16{}
			return nil
		}

		parts := strings.Split(str, ",")
		result := make([]int16, len(parts))
		for i, part := range parts {
			var val int16
			_, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &val)
			if err != nil {
				return err
			}
			result[i] = val
		}
		*a = result
		return nil
	case []int16:
		// Handle native slice
		*a = v
		return nil
	case []int64:
		// Handle int64 slice (common conversion)
		result := make([]int16, len(v))
		for i, val := range v {
			result[i] = int16(val)
		}
		*a = result
		return nil
	default:
		return fmt.Errorf("incompatible type for Int16Array: %T", v)
	}
}

func (a Int16Array) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}

	strs := make([]string, len(a))
	for i, v := range a {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ",")), nil
}

type Measurements struct {
	Bust      float32 `json:"bust"`
	UnderBust float32 `json:"underBust"`
	Waist     float32 `json:"waist"`
	Hip       float32 `json:"hip"`
}

type Preferences struct {
	gorm.Model
	Style            Int16Array   `json:"-" gorm:"type:smallint[]"`
	FavoriteColors   Int16Array `json:"-" gorm:"type:smallint[]"`
	PreferredBra     int16      `json:"-" gorm:"default:1"`
	PreferredModel   int16      `json:"-" gorm:"default:1"`
	PreferredPanties int16      `json:"-" gorm:"default:1"`
	Size             int16      `json:"-" gorm:"default:1"`
	AllowedModels    Int16Array   `json:"allowedModels" gorm:"type:smallint[]"`
	NotAllowedModels string       `json:"notAllowedModels" gorm:"type:text"`
	Notes            string       `json:"notes" gorm:"type:text"`
	Measurements Measurements `json:"measurements" gorm:"embedded;embeddedPrefix:measurements_"`
}

// MarshalJSON implements custom JSON marshaling for Preferences
func (p Preferences) MarshalJSON() ([]byte, error) {
	// Convert style IDs to names
	var styleNames []string
	if len(p.Style) > 0 {
		for _, styleID := range p.Style {
			if styleName, ok := enums.StyleNames[int(styleID)]; ok {
				styleNames = append(styleNames, styleName)
			}
		}
	}

	// Convert favorite color IDs to names
	var favoriteColorNames []string
	if len(p.FavoriteColors) > 0 {
		for _, colorID := range p.FavoriteColors {
			if colorName, ok := enums.ColorNames[colorID]; ok {
				favoriteColorNames = append(favoriteColorNames, colorName)
			}
		}
	}

	// Convert allowed model IDs to names
	var allowedModelNames []string
	if len(p.AllowedModels) > 0 {
		for _, modelID := range p.AllowedModels {
			if modelName, ok := preferred.ModelNames[modelID]; ok {
				allowedModelNames = append(allowedModelNames, modelName)
			}
		}
	}

	type Alias Preferences
	return json.Marshal(&struct {
		*Alias
		Style            []string `json:"style"`
		FavoriteColors   []string `json:"favoriteColors"`
		PreferredBra     string   `json:"preferredBra"`
		PreferredModel   string   `json:"preferredModel"`
		PreferredPanties string   `json:"preferredPanties"`
		Size             string   `json:"size"`
		AllowedModels    []string `json:"allowedModels"`
	}{
		Alias:            (*Alias)(&p),
		Style:            styleNames,
		FavoriteColors:   favoriteColorNames,
		PreferredBra:     getBraName(p.PreferredBra),
		PreferredModel:   getModelName(p.PreferredModel),
		PreferredPanties: getPantieName(p.PreferredPanties),
		Size:             getSizeName(p.Size),
		AllowedModels:    allowedModelNames,
	})
}

func getBraName(id int16) string {
	if name, ok := preferred.BraNames[int(id)]; ok {
		return name
	}
	return "Unknown"
}

func getModelName(id int16) string {
	if name, ok := preferred.ModelNames[id]; ok {
		return name
	}
	return "Unknown"
}

func getPantieName(id int16) string {
	if name, ok := preferred.PantieNames[int(id)]; ok {
		return name
	}
	return "Unknown"
}

func getSizeName(id int16) string {
	if name, ok := enums.SizeNames[id]; ok {
		return name
	}
	return "Unknown"
}
