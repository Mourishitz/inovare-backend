package repositories

import (
	"errors"

	"inovare-backend/database"
	"inovare-backend/models"
	"inovare-backend/requests"
	"inovare-backend/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user requests.CreateUserRequest) (*models.User, error)
	Update(user models.User) (*models.User, error)
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

// Create implements UserRepository.
func (u *userRepository) Create(user requests.CreateUserRequest) (*models.User, error) {
	newUser := models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	if err := u.db.Create(&newUser).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, utils.ErrDuplicateEmail
			}
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, utils.ErrDuplicateEmail
		}

		return nil, err
	}

	return &newUser, nil
}

// Delete implements UserRepository.
func (u *userRepository) Delete(id int) error {
	panic("unimplemented")
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

// GetByID implements UserRepository.
func (u *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User

	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

// Update implements UserRepository.
func (u *userRepository) Update(user models.User) (*models.User, error) {
	panic("unimplemented")
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}
