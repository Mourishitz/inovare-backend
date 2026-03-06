package repositories

import (
	"errors"

	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByID(id int) (*models.Product, error)
	GetAll() ([]models.Product, error)
	GetAllPaginated(page, pageSize int) ([]models.Product, int64, error)
	Search(query string, catalogID *uint) ([]models.Product, error)
	Create(product requests.CreateProductRequest) (*models.Product, error)
	Update(id int, updates requests.UpdateProductRequest) (*models.Product, error)
	Delete(id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		db: database.DB,
	}
}

// GetByID implements ProductRepository.
func (r *productRepository) GetByID(id int) (*models.Product, error) {
	var product models.Product

	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

// GetAll implements ProductRepository.
func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// GetAllPaginated implements ProductRepository.
func (r *productRepository) GetAllPaginated(page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count total records
	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated results
	if err := r.db.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Search implements ProductRepository.
func (r *productRepository) Search(query string, catalogID *uint) ([]models.Product, error) {
	var products []models.Product

	db := r.db.Where("(name ILIKE ? OR description ILIKE ?)", "%"+query+"%", "%"+query+"%")

	if catalogID != nil {
		db = db.Where("is_exclusive = false OR (is_exclusive = true AND catalog_id = ?)", *catalogID)
	} else {
		db = db.Where("is_exclusive = false")
	}

	if err := db.Limit(5).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// Create implements ProductRepository.
func (r *productRepository) Create(product requests.CreateProductRequest) (*models.Product, error) {
	newProduct := models.Product{
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		IsExclusive: product.IsExclusive,
		CatalogID:   product.CatalogID,
	}

	if err := r.db.Create(&newProduct).Error; err != nil {
		return nil, err
	}

	return &newProduct, nil
}

// Update implements ProductRepository.
func (r *productRepository) Update(id int, updates requests.UpdateProductRequest) (*models.Product, error) {
	product, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	updateData := make(map[string]interface{})

	if updates.Name != nil {
		updateData["name"] = *updates.Name
	}
	if updates.Description != nil {
		updateData["description"] = *updates.Description
	}
	if updates.ImageURL != nil {
		updateData["image_url"] = *updates.ImageURL
	}
	if updates.IsExclusive != nil {
		updateData["is_exclusive"] = *updates.IsExclusive
	}

	if len(updateData) > 0 {
		if err := r.db.Model(product).Updates(updateData).Error; err != nil {
			return nil, err
		}
	}

	return r.GetByID(id)
}

// Delete implements ProductRepository.
func (r *productRepository) Delete(id int) error {
	product, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(product).Error; err != nil {
		return err
	}

	return nil
}
