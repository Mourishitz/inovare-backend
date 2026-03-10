package services

import (
	cryptoRand "crypto/rand"
	"time"

	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
	"inovare-backend/utils"
)

type ShowerService interface {
	GetByID(id int) (*models.Shower, error)
	GetAll() ([]models.Shower, error)
	GetAllPaginated(page, pageSize int) ([]models.Shower, int64, error)
	GetByHostID(hostID uint) ([]models.Shower, error)
	Create(shower requests.CreateShowerRequest) (*models.Shower, error)
	CreateWithHost(shower requests.CreateShowerRequest, hostID uint) (*models.Shower, error)
	Update(id int, updates requests.UpdateShowerRequest) (*models.Shower, error)
	AddCatalog(showerID int, catalog requests.AddCatalogRequest) (*models.Shower, error)
	AddPreferences(showerID int, preferences requests.AddPreferencesRequest) (*models.Shower, error)
	GetCatalogWithProducts(showerID int) (*models.Catalog, []models.CatalogProduct, error)
	GetDashboardStats() (map[string]interface{}, error)
}

type showerService struct {
	showerRepo         repositories.ShowerRepository
	userRepo           repositories.UserRepository
	catalogRepo        repositories.CatalogRepository
	catalogProductRepo repositories.CatalogProductRepository
}

func NewShowerService() ShowerService {
	return &showerService{
		showerRepo:         repositories.NewShowerRepository(),
		userRepo:           repositories.NewUserRepository(),
		catalogRepo:        repositories.NewCatalogRepository(),
		catalogProductRepo: repositories.NewCatalogProductRepository(),
	}
}

// GetByID implements ShowerService.
func (s *showerService) GetByID(id int) (*models.Shower, error) {
	return s.showerRepo.GetByID(id)
}

// GetAll implements ShowerService.
func (s *showerService) GetAll() ([]models.Shower, error) {
	return s.showerRepo.GetAll()
}

// GetAllPaginated implements ShowerService.
func (s *showerService) GetAllPaginated(page, pageSize int) ([]models.Shower, int64, error) {
	return s.showerRepo.GetAllPaginated(page, pageSize)
}

// GetByHostID implements ShowerService.
func (s *showerService) GetByHostID(hostID uint) ([]models.Shower, error) {
	return s.showerRepo.GetByHostID(hostID)
}

// Create implements ShowerService.
func (s *showerService) Create(shower requests.CreateShowerRequest) (*models.Shower, error) {
	// Validate host exists
	_, err := s.userRepo.GetByID(int(shower.HostID))
	if err != nil {
		if err == utils.ErrUserNotFound {
			return nil, err
		}
		return nil, err
	}

	return s.showerRepo.Create(shower)
}

// CreateWithHost implements ShowerService with explicit host ID.
func (s *showerService) CreateWithHost(shower requests.CreateShowerRequest, hostID uint) (*models.Shower, error) {
	// Validate host exists
	_, err := s.userRepo.GetByID(int(hostID))
	if err != nil {
		if err == utils.ErrUserNotFound {
			return nil, err
		}
		return nil, err
	}

	// Create request with host ID
	createReq := requests.CreateShowerRequest{
		Guests:      shower.Guests,
		ShowerDate:  shower.ShowerDate,
		WeddingDate: shower.WeddingDate,
		Location:    shower.Location,
		HostID:      hostID,
	}

	return s.showerRepo.Create(createReq)
}

// Update implements ShowerService.
func (s *showerService) Update(id int, updates requests.UpdateShowerRequest) (*models.Shower, error) {
	// Validate shower exists
	_, err := s.showerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.showerRepo.Update(id, updates)
}

// AddCatalog implements ShowerService.
func (s *showerService) AddCatalog(showerID int, catalogReq requests.AddCatalogRequest) (*models.Shower, error) {
	// Validate shower exists
	_, err := s.showerRepo.GetByID(showerID)
	if err != nil {
		return nil, err
	}

	// Generate unique ID for catalog - retry until unique
	var uniqueID string
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		uniqueID = generateUniqueID()
		exists, err := s.catalogRepo.ExistsByURL(uniqueID)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		// If all retries exhausted, increase length
		if i == maxRetries-1 {
			uniqueID = generateRandomString(24) // Increase to 24 chars
		}
	}

	catalog := &models.Catalog{
		URL:      uniqueID, // Store only the unique ID
		Package:  catalogReq.Package,
		Approved: false,
	}

	if err := s.showerRepo.AddCatalog(showerID, catalog); err != nil {
		return nil, err
	}

	return s.showerRepo.GetByID(showerID)
}

// generateUniqueID creates a random unique identifier
func generateUniqueID() string {
	return generateRandomString(16)
}

// generateRandomString generates a random alphanumeric string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[randomInt(len(charset))]
	}
	return string(result)
}

// randomInt generates a random integer between 0 and max-1
func randomInt(max int) int {
	// Use crypto/rand for better randomness
	var b [1]byte
	_, err := cryptoRand.Read(b[:])
	if err != nil {
		// Fallback to time-based seed if crypto/rand fails
		return int(time.Now().UnixNano()) % max
	}
	return int(b[0]) % max
}

// AddPreferences implements ShowerService.
func (s *showerService) AddPreferences(showerID int, preferencesReq requests.AddPreferencesRequest) (*models.Shower, error) {
	// Validate shower exists
	_, err := s.showerRepo.GetByID(showerID)
	if err != nil {
		return nil, err
	}

	preferences := &models.Preferences{
		Style:            models.Int16Array(preferencesReq.Style),
		FavoriteColors:   models.Int16Array(preferencesReq.FavoriteColors),
		PreferredBra:     preferencesReq.PreferredBra,
		PreferredModel:   preferencesReq.PreferredModel,
		PreferredPanties: preferencesReq.PreferredPanties,
		Size:             preferencesReq.Size,
		AllowedModels:    models.Int16Array(preferencesReq.AllowedModels),
		NotAllowedModels: preferencesReq.NotAllowedModels,
		Notes:            preferencesReq.Notes,
		Measurements: models.Measurements{
			Bust:      preferencesReq.Measurements.Bust,
			UnderBust: preferencesReq.Measurements.UnderBust,
			Waist:     preferencesReq.Measurements.Waist,
			Hip:       preferencesReq.Measurements.Hip,
		},
	}

	if err := s.showerRepo.AddPreferences(showerID, preferences); err != nil {
		return nil, err
	}

	return s.showerRepo.GetByID(showerID)
}

// GetCatalogWithProducts retrieves the catalog and all its products for a shower
func (s *showerService) GetCatalogWithProducts(showerID int) (*models.Catalog, []models.CatalogProduct, error) {
	// Validate shower exists
	shower, err := s.showerRepo.GetByID(showerID)
	if err != nil {
		return nil, nil, err
	}

	// Check if shower has a catalog
	if shower.CatalogID == nil {
		return nil, nil, utils.ErrCatalogNotFound
	}

	// Get catalog
	catalog, err := s.catalogRepo.GetByID(int(*shower.CatalogID))
	if err != nil {
		return nil, nil, err
	}

	// Get all products in the catalog
	products, err := s.catalogProductRepo.GetByCatalogID(int(*shower.CatalogID))
	if err != nil {
		return nil, nil, err
	}

	return catalog, products, nil
}

// GetDashboardStats implements ShowerService.
func (s *showerService) GetDashboardStats() (map[string]interface{}, error) {
	totalShowers, approvedCatalogs, notApprovedCatalogs, recentShowers, err := s.showerRepo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_showers":           totalShowers,
		"approved_catalogs":       approvedCatalogs,
		"not_approved_catalogs":   notApprovedCatalogs,
		"recent_showers":          recentShowers,
	}, nil
}
