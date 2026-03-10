package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogService interface {
	GetByID(id int) (*models.Catalog, error)
	GetProductsByURL(url string) (*models.Catalog, []models.CatalogProduct, error)
	Approve(id int) (*models.Catalog, error)
	RegisterChanges(id int) (*models.Catalog, error)
}

type catalogService struct {
	catalogRepo        repositories.CatalogRepository
	catalogProductRepo repositories.CatalogProductRepository
	showerRepo         repositories.ShowerRepository
	emailService       EmailService
}

func NewCatalogService() CatalogService {
	return &catalogService{
		catalogRepo:        repositories.NewCatalogRepository(),
		catalogProductRepo: repositories.NewCatalogProductRepository(),
		showerRepo:         repositories.NewShowerRepository(),
		emailService:       NewEmailService(),
	}
}

// GetByID returns a catalog by its ID
func (s *catalogService) GetByID(id int) (*models.Catalog, error) {
	catalog, err := s.catalogRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}
	return catalog, nil
}

// GetProductsByURL fetches a catalog by URL and returns its products if approved
func (s *catalogService) GetProductsByURL(url string) (*models.Catalog, []models.CatalogProduct, error) {
	catalog, err := s.catalogRepo.GetByURL(url)
	if err != nil {
		return nil, nil, err
	}

	if !catalog.Approved {
		return nil, nil, utils.ErrCatalogNotApproved
	}

	products, err := s.catalogProductRepo.GetByCatalogID(int(catalog.ID))
	if err != nil {
		return nil, nil, err
	}

	return catalog, products, nil
}

// Approve approves a catalog
func (s *catalogService) Approve(id int) (*models.Catalog, error) {
	catalog, err := s.catalogRepo.Approve(id)
	if err != nil {
		if err == utils.ErrCatalogNotFound || err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}
	return catalog, nil
}

// RegisterChanges bumps updated_at and notifies the catalog host via e-mail
func (s *catalogService) RegisterChanges(id int) (*models.Catalog, error) {
	catalog, err := s.catalogRepo.TouchUpdatedAt(id)
	if err != nil {
		return nil, err
	}

	shower, err := s.showerRepo.GetByCatalogID(catalog.ID)
	if err != nil {
		// Non-fatal: catalog exists but may not be attached to a shower yet
		return catalog, nil
	}

	_ = s.emailService.SendCatalogChangesNotification(shower.Host.Email, shower.Host.Username, catalog.ID)

	return catalog, nil
}
