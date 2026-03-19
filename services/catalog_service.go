package services

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

type CatalogService interface {
	GetByID(id int) (*models.Catalog, error)
	GetProductsByURL(url string) (*models.Catalog, []models.CatalogProduct, error)
	CreatePublicPurchase(ctx context.Context, slug string, req requests.PublicCatalogPurchaseRequest) (*requests.PublicCatalogPurchaseResponse, error)
	Approve(id int, userID int) (*models.Catalog, error)
	RegisterChanges(id int) (*models.Catalog, error)
}

type catalogService struct {
	catalogRepo        repositories.CatalogRepository
	catalogProductRepo repositories.CatalogProductRepository
	showerRepo         repositories.ShowerRepository
	emailService       EmailService
	paymentProvider    PaymentProviderService
	now                func() time.Time
}

func NewCatalogService() CatalogService {
	return &catalogService{
		catalogRepo:        repositories.NewCatalogRepository(),
		catalogProductRepo: repositories.NewCatalogProductRepository(),
		showerRepo:         repositories.NewShowerRepository(),
		emailService:       NewEmailService(),
		paymentProvider:    NewAbacatePayService(),
		now:                time.Now,
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

// CreatePublicPurchase creates an Abacate Pay PIX checkout for a public catalog product.
func (s *catalogService) CreatePublicPurchase(ctx context.Context, slug string, req requests.PublicCatalogPurchaseRequest) (*requests.PublicCatalogPurchaseResponse, error) {
	productID, err := parsePublicCatalogProductID(req.ProductID)
	if err != nil {
		return nil, err
	}

	if err := validatePurchaseCustomer(req.Customer); err != nil {
		return nil, err
	}

	if err := validateCurrentURL(req.CurrentURL); err != nil {
		return nil, err
	}

	catalog, products, err := s.GetProductsByURL(slug)
	if err != nil {
		return nil, err
	}

	catalogProduct, err := findCatalogProduct(products, productID)
	if err != nil {
		return nil, err
	}

	if catalogProduct.IsBought {
		return nil, utils.ErrCatalogProductAlreadyBought
	}

	checkoutURL, err := s.paymentProvider.CreatePIXBilling(ctx, CreatePIXBillingRequest{
		ExternalID:    fmt.Sprintf("catalog:%d:product:%d:checkout:%d", catalog.ID, catalogProduct.ProductID, s.now().UnixMilli()),
		Name:          catalogProduct.Product.Name,
		Description:   catalogProduct.Product.Description,
		PriceInCents:  int(math.Round(catalogProduct.Price)),
		ReturnURL:     req.CurrentURL,
		CompletionURL: req.CurrentURL,
		Customer:      req.Customer,
	})
	if err != nil {
		return nil, err
	}

	return &requests.PublicCatalogPurchaseResponse{
		CheckoutURL: checkoutURL,
	}, nil
}

// Approve approves a catalog if the authenticated user owns the related shower.
func (s *catalogService) Approve(id int, userID int) (*models.Catalog, error) {
	catalog, err := s.catalogRepo.GetByID(id)
	if err != nil {
		if err == utils.ErrCatalogNotFound || err == gorm.ErrRecordNotFound {
			return nil, utils.ErrCatalogNotFound
		}
		return nil, err
	}

	shower, err := s.showerRepo.GetByCatalogID(catalog.ID)
	if err != nil {
		if err == utils.ErrShowerNotFound || err == gorm.ErrRecordNotFound {
			return nil, utils.ErrShowerNotFound
		}
		return nil, err
	}

	if shower.HostID != uint(userID) {
		return nil, utils.ErrUnauthorizedShowerAccess
	}

	catalog, err = s.catalogRepo.Approve(id)
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

func parsePublicCatalogProductID(value string) (uint, error) {
	productID, err := strconv.ParseUint(strings.TrimSpace(value), 10, 32)
	if err != nil || productID == 0 {
		return 0, utils.ErrInvalidProductID
	}

	return uint(productID), nil
}

func validateCurrentURL(rawURL string) error {
	parsedURL, err := url.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil || parsedURL == nil || parsedURL.Host == "" {
		return utils.ErrInvalidCurrentURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return utils.ErrInvalidCurrentURL
	}

	return nil
}

func validatePurchaseCustomer(customer requests.PublicCatalogPurchaseCustomer) error {
	if strings.TrimSpace(customer.Name) == "" ||
		strings.TrimSpace(customer.Email) == "" ||
		strings.TrimSpace(customer.Cellphone) == "" ||
		strings.TrimSpace(customer.TaxID) == "" {
		return utils.ErrInvalidCustomerData
	}

	return nil
}

func findCatalogProduct(products []models.CatalogProduct, productID uint) (*models.CatalogProduct, error) {
	for i := range products {
		if products[i].ProductID == productID {
			return &products[i], nil
		}
	}

	return nil, utils.ErrCatalogProductNotFound
}
