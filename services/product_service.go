package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
)

type ProductService interface {
	GetByID(id int) (*models.Product, error)
	GetAll() ([]models.Product, error)
	GetAllPaginated(page, pageSize int) ([]models.Product, int64, error)
	Search(query string, catalogID *uint) ([]models.Product, error)
	Create(product requests.CreateProductRequest) (*models.Product, error)
	Update(id int, updates requests.UpdateProductRequest) (*models.Product, error)
	Delete(id int) error
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService() ProductService {
	return &productService{
		productRepo: repositories.NewProductRepository(),
	}
}

// GetByID implements ProductService.
func (s *productService) GetByID(id int) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

// GetAll implements ProductService.
func (s *productService) GetAll() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

// GetAllPaginated implements ProductService.
func (s *productService) GetAllPaginated(page, pageSize int) ([]models.Product, int64, error) {
	return s.productRepo.GetAllPaginated(page, pageSize)
}

// Search implements ProductService.
func (s *productService) Search(query string, catalogID *uint) ([]models.Product, error) {
	return s.productRepo.Search(query, catalogID)
}

// Create implements ProductService.
func (s *productService) Create(product requests.CreateProductRequest) (*models.Product, error) {
	return s.productRepo.Create(product)
}

// Update implements ProductService.
func (s *productService) Update(id int, updates requests.UpdateProductRequest) (*models.Product, error) {
	return s.productRepo.Update(id, updates)
}

// Delete implements ProductService.
func (s *productService) Delete(id int) error {
	return s.productRepo.Delete(id)
}
