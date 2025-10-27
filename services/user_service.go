package services

import (
	"inovare-backend/models"
	"inovare-backend/repositories"
	"inovare-backend/requests"
)

type UserService interface {
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user requests.CreateUserRequest) (*models.User, error)
	Update(user models.User) (*models.User, error)
	Delete(id int) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *userService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Create(user requests.CreateUserRequest) (*models.User, error) {
	newUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userService) Update(user models.User) (*models.User, error) {
	panic("unimplemented")
}

func (s *userService) Delete(id int) error {
	panic("unimplemented")
}
