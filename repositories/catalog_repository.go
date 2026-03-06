package repositories

import (
	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	ExistsByURL(url string) (bool, error)
	GetByID(id int) (*models.Catalog, error)
	Approve(id int) (*models.Catalog, error)
	TouchUpdatedAt(id int) (*models.Catalog, error)
}

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository() CatalogRepository {
	return &catalogRepository{
		db: database.DB,
	}
}

// ExistsByURL checks if a catalog with the given URL already exists
func (r *catalogRepository) ExistsByURL(url string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Catalog{}).Where("url = ?", url).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByID retrieves a catalog by ID
func (r *catalogRepository) GetByID(id int) (*models.Catalog, error) {
	var catalog models.Catalog
	err := r.db.First(&catalog, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &catalog, nil
}

// Approve sets the approved status of a catalog to true
func (r *catalogRepository) Approve(id int) (*models.Catalog, error) {
	catalog, err := r.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	catalog.Approved = true
	if err := r.db.Save(catalog).Error; err != nil {
		return nil, err
	}

	return catalog, nil
}

// TouchUpdatedAt bumps the updated_at timestamp of a catalog to now
func (r *catalogRepository) TouchUpdatedAt(id int) (*models.Catalog, error) {
	catalog, err := r.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	if err := r.db.Model(catalog).UpdateColumn("updated_at", gorm.Expr("NOW()")).Error; err != nil {
		return nil, err
	}

	return r.GetByID(id)
}
