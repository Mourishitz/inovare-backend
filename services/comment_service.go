package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/utils"
)

type CommentService interface {
	AddComment(catalogID int, authorID uint, content string) (*models.Comment, error)
	ListComments(catalogID int) ([]models.Comment, error)
}

type commentService struct {
	commentRepo repositories.CommentRepository
	catalogRepo repositories.CatalogRepository
	showerRepo  repositories.ShowerRepository
}

func NewCommentService() CommentService {
	return &commentService{
		commentRepo: repositories.NewCommentRepository(),
		catalogRepo: repositories.NewCatalogRepository(),
		showerRepo:  repositories.NewShowerRepository(),
	}
}

// AddComment creates a comment on a catalog. Only the shower host is allowed.
func (s *commentService) AddComment(catalogID int, authorID uint, content string) (*models.Comment, error) {
	catalog, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		return nil, utils.ErrCatalogNotFound
	}

	shower, err := s.showerRepo.GetByCatalogID(catalog.ID)
	if err != nil {
		return nil, utils.ErrShowerNotFound
	}

	if shower.HostID != authorID {
		return nil, utils.ErrUnauthorizedShowerAccess
	}

	return s.commentRepo.Create(catalogID, authorID, content)
}

// ListComments returns all comments for a catalog.
func (s *commentService) ListComments(catalogID int) ([]models.Comment, error) {
	_, err := s.catalogRepo.GetByID(catalogID)
	if err != nil {
		return nil, utils.ErrCatalogNotFound
	}
	return s.commentRepo.GetByCatalogID(catalogID)
}
