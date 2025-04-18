package services

import (
	"context"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/utils"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int) error
	AddUser(ctx context.Context, user *models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ModifyUser(ctx context.Context, user *models.User) error
}

type userService struct {
	db UserRepository
}

func NewUserService(db UserRepository) *userService {
	return &userService{db: db}
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.db.DeleteUser(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, email string, password string, admin bool) (*models.User, error) {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: hashPassword,
		Admin:    admin,
	}

	id, err := s.db.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Id = id
	return user, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return s.db.GetUser(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.db.GetUserByEmail(ctx, email)
}

func (s *userService) GetAllUsers(ctx context.Context) (users []models.User, err error) {
	return s.db.GetAllUsers(ctx)
}

func (s *userService) ModifyUser(ctx context.Context, user *models.User) error {
	return s.db.ModifyUser(ctx, user)

}
