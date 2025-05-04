package services

import (
	"context"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/repositories"
	"ing-soft-2-tp1/internal/utils"
)

type UserService interface {
	DeleteUser(ctx context.Context, id int) error
	CreateUser(ctx context.Context, request models.CreateUserRequest, admin bool) (*models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context) (users []models.User, err error)
	ModifyUser(ctx context.Context, user *models.User) error
	BlockUser(ctx context.Context, id int) error
	ModifyLocation(ctx context.Context, id int, newLocation string) error
	Login(ctx context.Context, email string, password string) (*models.User, error)
	IsUserBlocked(ctx context.Context, id int) (bool, error)
}

type userService struct {
	userRepo        repositories.UserRepository
	blockedUserRepo repositories.BlockedUserRepository
}

func NewUserService(db repositories.UserRepository) *userService {
	return &userService{userRepo: db}
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, request models.CreateUserRequest, admin bool) (*models.User, error) {
	hashPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    request.Email,
		Password: hashPassword,
		Name:     request.Name,
		Surname:  request.Surname,
		Admin:    admin,
	}

	id, err := s.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		user.Id = id
		return user, nil
	}

	return createdUser, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetUser(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *userService) GetAllUsers(ctx context.Context) (users []models.User, err error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *userService) ModifyUser(ctx context.Context, user *models.User) error {
	return s.userRepo.ModifyUser(ctx, user)
}

func (s *userService) ModifyLocation(ctx context.Context, id int, newLocation string) error {
	return s.userRepo.ModifyLocation(ctx, id, newLocation)
}

// Login handles user login
// It checks if the user exists and if the password is correct
func (s *userService) Login(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := utils.CompareHashPassword(user.Password, password); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) IsUserBlocked(ctx context.Context, id int) (bool, error) {
	block, _, err := s.blockedUserRepo.IsUserBlocked(ctx, id)
	return block, err
}
