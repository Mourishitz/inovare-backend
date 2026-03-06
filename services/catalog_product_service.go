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
	CreateExclusiveProduct(catalogID int, req requests.CreateExclusiveProductRequest) (*models.CatalogProduct, error)
	GetCatalogIDByProductID(productID uint) (*uint, error)
	ListCatalogProducts(catalogID int) ([]models.CatalogProduct, error)
	UpdateCatalogProduct(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error)
	DetachProduct(catalogID, productID int) error
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
	product, err := s.productRepo.GetByID(int(req.ProductID))
	if err != nil {
		if err == utils.ErrProductNotFound {
			return nil, utils.ErrProductNotFound
		}
		return nil, err
	}

	// Reject if product is exclusive and already assigned to any catalog
	if product.IsExclusive {
		exists, err := s.catalogProductRepo.ProductExistsInAnyCatalog(product.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, utils.ErrProductIsExclusive
		}
	}

	// Attach product to catalog
	return s.catalogProductRepo.AttachProduct(catalogID, req)
}

// CreateExclusiveProduct creates a product marked as exclusive and attaches it to a catalog
func (s *catalogProductService) CreateExclusiveProduct(catalogID int, req requests.CreateExclusiveProductRequest) (*models.CatalogProduct, error) {
	// Validate catalog exists
	_, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	// Create product with exclusive flag
	catalogIDUint := uint(catalogID)
	product, err := s.productRepo.Create(requests.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		IsExclusive: true,
		CatalogID:   &catalogIDUint,
	})
	if err != nil {
		return nil, err
	}

	// Attach the newly created exclusive product to the catalog
	return s.catalogProductRepo.AttachProduct(catalogID, requests.AttachProductToCatalogRequest{
		ProductID: product.ID,
		Price:     req.Price,
		IsBought:  req.IsBought,
	})
}

// GetCatalogIDByProductID returns the catalog ID a product is attached to, or nil if none
func (s *catalogProductService) GetCatalogIDByProductID(productID uint) (*uint, error) {
	return s.catalogProductRepo.GetCatalogIDByProductID(productID)
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
func (s *catalogProductService) DetachProduct(catalogID, productID int) error {
	return s.catalogProductRepo.DeleteByCatalogAndProductID(catalogID, productID)
}
