package services

import (
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/utils"
)

type UserRepository interface {
	GetUser(id int) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(id int) error
	AddUser(user *models.User) (int, error)
	GetUserByEmail(email string) (*models.User, error)
	ModifyUser(user *models.User) error
}

type userService struct {
	db UserRepository
}

func NewUserService(db UserRepository) *userService {
	return &userService{db: db}
}

func (s *userService) DeleteUser(id int) error {
	return s.db.DeleteUser(id)
}

func (s *userService) CreateUser(email string, password string, admin bool) (*models.User, error) {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: hashPassword,
		Admin:    admin,
	}

	id, err := s.db.AddUser(user)
	if err != nil {
		return nil, err
	}

	user.Id = id
	return user, nil
}

func (s *userService) GetUserById(id int) (*models.User, error) {
	return s.db.GetUser(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.db.GetUserByEmail(email)
}

func (s *userService) GetAllUsers() (users []models.User, err error) {
	return s.db.GetAllUsers()
}

func (s *userService) ModifyUser(user *models.User) error {
	return s.db.ModifyUser(user)

}
