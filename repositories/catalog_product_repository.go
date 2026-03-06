package repositories

import (
	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogProductRepository interface {
	AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error)
	ProductExistsInCatalog(catalogID int, productID uint) (bool, error)
	ProductExistsInAnyCatalog(productID uint) (bool, error)
	GetCatalogIDByProductID(productID uint) (*uint, error)
	GetByID(id int) (*models.CatalogProduct, error)
	GetByCatalogID(catalogID int) ([]models.CatalogProduct, error)
	Update(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error)
	Delete(id int) error
	DeleteByCatalogAndProductID(catalogID, productID int) error
}

type catalogProductRepository struct {
	db *gorm.DB
}

func NewCatalogProductRepository() CatalogProductRepository {
	return &catalogProductRepository{
		db: database.DB,
	}
}

// AttachProduct attaches a product to a catalog
func (r *catalogProductRepository) AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error) {
	// Check if product already exists in this catalog
	exists, err := r.ProductExistsInCatalog(catalogID, req.ProductID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrProductAlreadyInCatalog
	}

	catalogProduct := &models.CatalogProduct{
		CatalogID: uint(catalogID),
		ProductID: req.ProductID,
		Price:     req.Price,
		IsBought:  req.IsBought,
	}

	if err := r.db.Create(catalogProduct).Error; err != nil {
		return nil, err
	}

	// Reload with relations
	if err := r.db.Preload("Product").Preload("Catalog").First(catalogProduct, catalogProduct.ID).Error; err != nil {
		return nil, err
	}

	return catalogProduct, nil
}

// GetCatalogIDByProductID returns the catalog ID a product is attached to, or nil if none
func (r *catalogProductRepository) GetCatalogIDByProductID(productID uint) (*uint, error) {
	var catalogProduct models.CatalogProduct
	err := r.db.Select("catalog_id").Where("product_id = ?", productID).First(&catalogProduct).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &catalogProduct.CatalogID, nil
}

// ProductExistsInAnyCatalog checks if an exclusive product is already assigned to a catalog
func (r *catalogProductRepository) ProductExistsInAnyCatalog(productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).
		Where("id = ? AND catalog_id IS NOT NULL", productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ProductExistsInCatalog checks if a product already exists in a catalog
func (r *catalogProductRepository) ProductExistsInCatalog(catalogID int, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.CatalogProduct{}).
		Where("catalog_id = ? AND product_id = ?", catalogID, productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetByID retrieves a catalog product by ID
func (r *catalogProductRepository) GetByID(id int) (*models.CatalogProduct, error) {
	var catalogProduct models.CatalogProduct
	err := r.db.Preload("Product").Preload("Catalog").First(&catalogProduct, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogProductNotFound
		}
		return nil, err
	}
	return &catalogProduct, nil
}

// GetByCatalogID retrieves all products in a catalog
func (r *catalogProductRepository) GetByCatalogID(catalogID int) ([]models.CatalogProduct, error) {
	var catalogProducts []models.CatalogProduct
	err := r.db.Preload("Product").Preload("Catalog").
		Where("catalog_id = ?", catalogID).
		Find(&catalogProducts).Error
	if err != nil {
		return nil, err
	}
	return catalogProducts, nil
}

// Update updates a catalog product
func (r *catalogProductRepository) Update(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error) {
	catalogProduct, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if updates.Price != nil {
		updateData["price"] = *updates.Price
	}
	if updates.IsBought != nil {
		updateData["is_bought"] = *updates.IsBought
	}

	if len(updateData) > 0 {
		if err := r.db.Model(catalogProduct).Updates(updateData).Error; err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

// Delete removes a catalog product (detaches product from catalog)
func (r *catalogProductRepository) Delete(id int) error {
	catalogProduct, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(catalogProduct).Error; err != nil {
		return err
	}

	return nil
}

// DeleteByCatalogAndProductID removes a product from a catalog by catalog_id + product_id
func (r *catalogProductRepository) DeleteByCatalogAndProductID(catalogID, productID int) error {
	var catalogProduct models.CatalogProduct
	err := r.db.Where("catalog_id = ? AND product_id = ?", catalogID, productID).First(&catalogProduct).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrCatalogProductNotFound
		}
		return err
	}
	return r.db.Delete(&catalogProduct).Error
}
