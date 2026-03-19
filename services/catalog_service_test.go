package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"inovare-backend/models"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"gorm.io/gorm"
)

func TestCatalogServiceCreatePublicPurchaseSuccess(t *testing.T) {
	paymentProvider := &fakePaymentProviderService{
		checkoutURL: "https://checkout.example.com/pix",
	}

	service := &catalogService{
		catalogRepo: &fakeCatalogRepository{
			catalogByURL: &models.Catalog{
				Model:    gorm.Model{ID: 1},
				Approved: true,
			},
		},
		catalogProductRepo: &fakeCatalogProductRepository{
			productsByCatalogID: []models.CatalogProduct{
				{
					Price:     5000,
					IsBought:  false,
					CatalogID: 1,
					ProductID: 2,
					Product: models.Product{
						Name:        "Gift",
						Description: "Special gift",
					},
				},
			},
		},
		showerRepo:      &fakeShowerRepository{},
		emailService:    &fakeEmailService{},
		paymentProvider: paymentProvider,
		now: func() time.Time {
			return time.UnixMilli(1773882016084)
		},
	}

	response, err := service.CreatePublicPurchase(context.Background(), "my-slug", requests.PublicCatalogPurchaseRequest{
		ProductID: "2",
		Customer: requests.PublicCatalogPurchaseCustomer{
			Name:      "Maria Silva",
			Email:     "maria@example.com",
			Cellphone: "5511999999999",
			TaxID:     "12345678901",
		},
		CurrentURL: "https://frontend.example.com/my-slug",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.CheckoutURL != "https://checkout.example.com/pix" {
		t.Fatalf("expected checkout URL to be returned, got %q", response.CheckoutURL)
	}

	if paymentProvider.lastRequest.ExternalID != "catalog:1:product:2:checkout:1773882016084" {
		t.Fatalf("expected external ID to include catalog/product/timestamp, got %q", paymentProvider.lastRequest.ExternalID)
	}

	if paymentProvider.lastRequest.PriceInCents != 5000 {
		t.Fatalf("expected price in cents to be 5000, got %d", paymentProvider.lastRequest.PriceInCents)
	}
}

func TestCatalogServiceCreatePublicPurchaseProductNotFound(t *testing.T) {
	service := &catalogService{
		catalogRepo: &fakeCatalogRepository{
			catalogByURL: &models.Catalog{Approved: true},
		},
		catalogProductRepo: &fakeCatalogProductRepository{
			productsByCatalogID: []models.CatalogProduct{},
		},
		showerRepo:      &fakeShowerRepository{},
		emailService:    &fakeEmailService{},
		paymentProvider: &fakePaymentProviderService{},
		now:             time.Now,
	}

	_, err := service.CreatePublicPurchase(context.Background(), "my-slug", requests.PublicCatalogPurchaseRequest{
		ProductID: "99",
		Customer: requests.PublicCatalogPurchaseCustomer{
			Name:      "Maria Silva",
			Email:     "maria@example.com",
			Cellphone: "5511999999999",
			TaxID:     "12345678901",
		},
		CurrentURL: "https://frontend.example.com/my-slug",
	})
	if !errors.Is(err, utils.ErrCatalogProductNotFound) {
		t.Fatalf("expected ErrCatalogProductNotFound, got %v", err)
	}
}

func TestCatalogServiceCreatePublicPurchaseAlreadyBought(t *testing.T) {
	service := &catalogService{
		catalogRepo: &fakeCatalogRepository{
			catalogByURL: &models.Catalog{Approved: true},
		},
		catalogProductRepo: &fakeCatalogProductRepository{
			productsByCatalogID: []models.CatalogProduct{
				{
					IsBought:  true,
					ProductID: 2,
				},
			},
		},
		showerRepo:      &fakeShowerRepository{},
		emailService:    &fakeEmailService{},
		paymentProvider: &fakePaymentProviderService{},
		now:             time.Now,
	}

	_, err := service.CreatePublicPurchase(context.Background(), "my-slug", requests.PublicCatalogPurchaseRequest{
		ProductID: "2",
		Customer: requests.PublicCatalogPurchaseCustomer{
			Name:      "Maria Silva",
			Email:     "maria@example.com",
			Cellphone: "5511999999999",
			TaxID:     "12345678901",
		},
		CurrentURL: "https://frontend.example.com/my-slug",
	})
	if !errors.Is(err, utils.ErrCatalogProductAlreadyBought) {
		t.Fatalf("expected ErrCatalogProductAlreadyBought, got %v", err)
	}
}

func TestCatalogServiceCreatePublicPurchaseInvalidInput(t *testing.T) {
	service := &catalogService{
		catalogRepo: &fakeCatalogRepository{
			catalogByURL: &models.Catalog{Approved: true},
		},
		catalogProductRepo: &fakeCatalogProductRepository{},
		showerRepo:         &fakeShowerRepository{},
		emailService:       &fakeEmailService{},
		paymentProvider:    &fakePaymentProviderService{},
		now:                time.Now,
	}

	_, err := service.CreatePublicPurchase(context.Background(), "my-slug", requests.PublicCatalogPurchaseRequest{
		ProductID: "2",
		Customer: requests.PublicCatalogPurchaseCustomer{
			Name:      "Maria Silva",
			Email:     "maria@example.com",
			Cellphone: "5511999999999",
			TaxID:     "12345678901",
		},
		CurrentURL: "/relative-url",
	})
	if !errors.Is(err, utils.ErrInvalidCurrentURL) {
		t.Fatalf("expected ErrInvalidCurrentURL, got %v", err)
	}
}

type fakeCatalogRepository struct {
	catalogByURL   *models.Catalog
	getByURLError  error
	catalogByID    *models.Catalog
	getByIDError   error
	approveCatalog *models.Catalog
	approveError   error
}

func (f *fakeCatalogRepository) ExistsByURL(url string) (bool, error) {
	return false, nil
}

func (f *fakeCatalogRepository) GetByID(id int) (*models.Catalog, error) {
	if f.getByIDError != nil {
		return nil, f.getByIDError
	}
	if f.catalogByID != nil {
		return f.catalogByID, nil
	}
	return &models.Catalog{}, nil
}

func (f *fakeCatalogRepository) GetByURL(url string) (*models.Catalog, error) {
	if f.getByURLError != nil {
		return nil, f.getByURLError
	}
	if f.catalogByURL != nil {
		return f.catalogByURL, nil
	}
	return nil, utils.ErrCatalogNotFound
}

func (f *fakeCatalogRepository) Approve(id int) (*models.Catalog, error) {
	if f.approveError != nil {
		return nil, f.approveError
	}
	if f.approveCatalog != nil {
		return f.approveCatalog, nil
	}
	return &models.Catalog{}, nil
}

func (f *fakeCatalogRepository) TouchUpdatedAt(id int) (*models.Catalog, error) {
	return &models.Catalog{}, nil
}

type fakeCatalogProductRepository struct {
	productsByCatalogID []models.CatalogProduct
	getByCatalogError   error
}

func (f *fakeCatalogProductRepository) AttachProduct(catalogID int, req requests.AttachProductToCatalogRequest) (*models.CatalogProduct, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) ProductExistsInCatalog(catalogID int, productID uint) (bool, error) {
	return false, nil
}

func (f *fakeCatalogProductRepository) ProductExistsInAnyCatalog(productID uint) (bool, error) {
	return false, nil
}

func (f *fakeCatalogProductRepository) GetCatalogIDByProductID(productID uint) (*uint, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) GetByID(id int) (*models.CatalogProduct, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) GetByCatalogID(catalogID int) ([]models.CatalogProduct, error) {
	if f.getByCatalogError != nil {
		return nil, f.getByCatalogError
	}
	return f.productsByCatalogID, nil
}

func (f *fakeCatalogProductRepository) GetByCatalogAndProductID(catalogID int, productID uint) (*models.CatalogProduct, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) Update(id int, updates requests.UpdateCatalogProductRequest) (*models.CatalogProduct, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) MarkAsBought(catalogID int, productID uint) (*models.CatalogProduct, error) {
	return nil, nil
}

func (f *fakeCatalogProductRepository) Delete(id int) error {
	return nil
}

func (f *fakeCatalogProductRepository) DeleteByCatalogAndProductID(catalogID, productID int) error {
	return nil
}

type fakeShowerRepository struct{}

func (f *fakeShowerRepository) GetByID(id int) (*models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) GetAll() ([]models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) GetAllPaginated(page, pageSize int) ([]models.Shower, int64, error) {
	return nil, 0, nil
}

func (f *fakeShowerRepository) GetByHostID(hostID uint) ([]models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) GetByCatalogID(catalogID uint) (*models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) Create(shower requests.CreateShowerRequest) (*models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) Update(id int, updates requests.UpdateShowerRequest) (*models.Shower, error) {
	return nil, nil
}

func (f *fakeShowerRepository) AddCatalog(showerID int, catalog *models.Catalog) error {
	return nil
}

func (f *fakeShowerRepository) AddPreferences(showerID int, preferences *models.Preferences) error {
	return nil
}

func (f *fakeShowerRepository) GetDashboardStats() (int64, int64, int64, []models.Shower, error) {
	return 0, 0, 0, nil, nil
}

type fakeEmailService struct{}

func (f *fakeEmailService) SendCatalogChangesNotification(hostEmail, hostName string, catalogID uint) error {
	return nil
}

type fakePaymentProviderService struct {
	checkoutURL string
	err         error
	lastRequest CreatePIXBillingRequest
}

func (f *fakePaymentProviderService) CreatePIXBilling(ctx context.Context, req CreatePIXBillingRequest) (string, error) {
	f.lastRequest = req
	if f.err != nil {
		return "", f.err
	}
	return f.checkoutURL, nil
}
