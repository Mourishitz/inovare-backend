package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogProductService interface {
	AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error)
	ListCatalogProducts(catalogID int) ([]models.CatalogProduct, error)
	UpdateCatalogProduct(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error)
	DetachProduct(id int) error
}

type catalogProductService struct {
	catalogProductRepo repositories.CatalogProductRepository
	catalogRepo        repositories.CatalogRepository
	productRepo        repositories.ProductRepository
}

func NewCatalogProductService() CatalogProductService {
	return &catalogProductService{
		catalogProductRepo: repositories.NewCatalogProductRepository(),
		catalogRepo:        repositories.NewCatalogRepository(),
		productRepo:        repositories.NewProductRepository(),
	}
}

// AttachProduct attaches a product to a catalog
func (s *catalogProductService) AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error) {
	// Validate catalog exists
	_, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	// Validate product exists
	_, err = s.productRepo.GetByID(int(req.ProductID))
	if err != nil {
		if err == utils.ErrProductNotFound {
			return nil, utils.ErrProductNotFound
		}
		return nil, err
	}

	// Attach product to catalog
	return s.catalogProductRepo.AttachProduct(catalogID, req)
}

// ListCatalogProducts lists all products in a catalog
func (s *catalogProductService) ListCatalogProducts(catalogID int) ([]models.CatalogProduct, error) {
	// Validate catalog exists
	_, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	return s.catalogProductRepo.GetByCatalogID(catalogID)
}

// UpdateCatalogProduct updates a catalog product
func (s *catalogProductService) UpdateCatalogProduct(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error) {
	return s.catalogProductRepo.Update(id, updates)
}

// DetachProduct detaches a product from a catalog
func (s *catalogProductService) DetachProduct(id int) error {
	return s.catalogProductRepo.Delete(id)
}
