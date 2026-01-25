package repositories

import (
	"errors"

	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type ShowerRepository interface {
	GetByID(id int) (*models.Shower, error)
	GetAll() ([]models.Shower, error)
	GetAllPaginated(page, pageSize int) ([]models.Shower, int64, error)
	Create(shower requests.CreateShowerRequest) (*models.Shower, error)
	Update(id int, updates requests.UpdateShowerRequest) (*models.Shower, error)
	AddCatalog(showerID int, catalog *models.Catalog) error
	AddPreferences(showerID int, preferences *models.Preferences) error
}

type showerRepository struct {
	db *gorm.DB
}

func NewShowerRepository() ShowerRepository {
	return &showerRepository{
		db: database.DB,
	}
}

// GetByID implements ShowerRepository.
func (r *showerRepository) GetByID(id int) (*models.Shower, error) {
	var shower models.Shower

	if err := r.db.Preload("Host").Preload("Catalog").Preload("Preferences").First(&shower, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrShowerNotFound
		}
		return nil, err
	}

	return &shower, nil
}

// GetAll implements ShowerRepository.
func (r *showerRepository) GetAll() ([]models.Shower, error) {
	var showers []models.Shower

	if err := r.db.Preload("Host").Preload("Catalog").Preload("Preferences").Find(&showers).Error; err != nil {
		return nil, err
	}

	return showers, nil
}

// GetAllPaginated implements ShowerRepository.
func (r *showerRepository) GetAllPaginated(page, pageSize int) ([]models.Shower, int64, error) {
	var showers []models.Shower
	var total int64

	// Count total records
	if err := r.db.Model(&models.Shower{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	if err := r.db.Preload("Host").Preload("Catalog").Preload("Preferences").
		Offset(offset).
		Limit(pageSize).
		Find(&showers).Error; err != nil {
		return nil, 0, err
	}

	return showers, total, nil
}

// Create implements ShowerRepository.
func (r *showerRepository) Create(shower requests.CreateShowerRequest) (*models.Shower, error) {
	newShower := models.Shower{
		Guests:      shower.Guests,
		ShowerDate:  shower.ShowerDate,
		WeddingDate: shower.WeddingDate,
		Location:    shower.Location,
		HostID:      shower.HostID,
	}

	if err := r.db.Create(&newShower).Error; err != nil {
		return nil, err
	}

	// Reload with relations
	return r.GetByID(int(newShower.ID))
}

// Update implements ShowerRepository.
func (r *showerRepository) Update(id int, updates requests.UpdateShowerRequest) (*models.Shower, error) {
	shower, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if updates.Guests != nil {
		updateData["guests"] = *updates.Guests
	}
	if updates.ShowerDate != nil {
		updateData["shower_date"] = *updates.ShowerDate
	}
	if updates.WeddingDate != nil {
		updateData["wedding_date"] = *updates.WeddingDate
	}
	if updates.Location != nil {
		updateData["location"] = *updates.Location
	}

	if len(updateData) > 0 {
		if err := r.db.Model(shower).Updates(updateData).Error; err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

// AddCatalog implements ShowerRepository.
func (r *showerRepository) AddCatalog(showerID int, catalog *models.Catalog) error {
	shower, err := r.GetByID(showerID)
	if err != nil {
		return err
	}

	if shower.CatalogID != nil {
		return utils.ErrCatalogAlreadyExists
	}

	if err := r.db.Create(catalog).Error; err != nil {
		return err
	}

	shower.CatalogID = &catalog.ID
	if err := r.db.Save(shower).Error; err != nil {
		return err
	}

	return nil
}

// AddPreferences implements ShowerRepository.
func (r *showerRepository) AddPreferences(showerID int, preferences *models.Preferences) error {
	shower, err := r.GetByID(showerID)
	if err != nil {
		return err
	}

	if shower.PreferencesID != nil {
		return utils.ErrPreferencesAlreadyExist
	}

	if err := r.db.Create(preferences).Error; err != nil {
		return err
	}

	shower.PreferencesID = &preferences.ID
	if err := r.db.Save(shower).Error; err != nil {
		return err
	}

	return nil
}
