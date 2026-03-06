package repositories

import (
	"inovare-backend/database"
	"inovare-backend/models"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(catalogID int, authorID uint, content string) (*models.Comment, error)
	GetByCatalogID(catalogID int) ([]models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository() CommentRepository {
	return &commentRepository{db: database.DB}
}

// Create adds a new comment to a catalog
func (r *commentRepository) Create(catalogID int, authorID uint, content string) (*models.Comment, error) {
	comment := models.Comment{
		Content:   content,
		AuthorID:  authorID,
		CatalogID: uint(catalogID),
	}
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	if err := r.db.Preload("Author").First(&comment, comment.ID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByCatalogID returns all comments for a catalog, ordered oldest first
func (r *commentRepository) GetByCatalogID(catalogID int) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("catalog_id = ?", catalogID).
		Preload("Author").
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
