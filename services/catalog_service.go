package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogService interface {
	Approve(id int) (*models.Catalog, error)
}

type catalogService struct {
	catalogRepo repositories.CatalogRepository
}

func NewCatalogService() CatalogService {
	return &catalogService{
		catalogRepo: repositories.NewCatalogRepository(),
	}
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
