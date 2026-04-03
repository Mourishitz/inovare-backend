package services

import (
	"fmt"

	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
	"inovare-backend/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type CatalogProductService interface {
	AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error)
	CreateExclusiveProduct(catalogID int, req requests.CreateExclusiveProductRequest) (*models.CatalogProduct, error)
	GetCatalogIDByProductID(productID uint) (*uint, error)
	ListCatalogProducts(catalogID int) ([]models.CatalogProduct, error)
	ListCatalogProductsWithFirstImage(catalogID int) ([]models.CatalogProduct, error)
	MarkAsBought(req requests.MarkCatalogProductAsBoughtRequest) (*models.CatalogProduct, error)
	MarkAsBoughtByExternalIDs(externalIDs []string) ([]models.CatalogProduct, error)
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
		Images:      req.Images,
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

func (s *catalogProductService) ListCatalogProductsWithFirstImage(catalogID int) ([]models.CatalogProduct, error) {
	_, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	return s.catalogProductRepo.GetByCatalogIDWithFirstImage(catalogID)
}

// MarkAsBought marks a catalog product as bought.
func (s *catalogProductService) MarkAsBought(req requests.MarkCatalogProductAsBoughtRequest) (*models.CatalogProduct, error) {
	return s.catalogProductRepo.MarkAsBought(int(req.CatalogID), req.ProductID)
}

// MarkAsBoughtByExternalIDs marks catalog products as bought using payment gateway external IDs.
func (s *catalogProductService) MarkAsBoughtByExternalIDs(externalIDs []string) ([]models.CatalogProduct, error) {
	updatedProducts := make([]models.CatalogProduct, 0, len(externalIDs))

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		catalogProductRepo := repositories.NewCatalogProductRepositoryWithDB(tx)

		for _, externalID := range externalIDs {
			catalogID, productID, err := parseCatalogProductExternalID(externalID)
			if err != nil {
				return err
			}

			catalogProduct, err := catalogProductRepo.MarkAsBought(int(catalogID), productID)
			if err != nil {
				return err
			}

			updatedProducts = append(updatedProducts, *catalogProduct)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProducts, nil
}

// UpdateCatalogProduct updates a catalog product
func (s *catalogProductService) UpdateCatalogProduct(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error) {
	return s.catalogProductRepo.Update(id, updates)
}

// DetachProduct detaches a product from a catalog
func (s *catalogProductService) DetachProduct(catalogID, productID int) error {
	return s.catalogProductRepo.DeleteByCatalogAndProductID(catalogID, productID)
}

func parseCatalogProductExternalID(externalID string) (uint, uint, error) {
	parts := strings.Split(externalID, ":")

	var catalogIDValue string
	var productIDValue string

	for i := 0; i < len(parts)-1; i++ {
		switch parts[i] {
		case "catalog":
			catalogIDValue = parts[i+1]
		case "product":
			productIDValue = parts[i+1]
		}
	}

	if catalogIDValue == "" || productIDValue == "" {
		return 0, 0, fmt.Errorf("%w: %s", utils.ErrInvalidWebhookExternalID, externalID)
	}

	catalogID, err := strconv.ParseUint(catalogIDValue, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %s", utils.ErrInvalidWebhookExternalID, externalID)
	}

	productID, err := strconv.ParseUint(productIDValue, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %s", utils.ErrInvalidWebhookExternalID, externalID)
	}

	return uint(catalogID), uint(productID), nil
}
