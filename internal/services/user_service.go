package services

import (
	"context"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"
)

type UserService interface {
	DeleteUser(ctx context.Context, id int) error
	CreateUser(ctx context.Context, request models.CreateUserRequest) (int, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context) (users []models.User, err error)
	ModifyUser(ctx context.Context, id int, user *models.User) error
	BlockUser(ctx context.Context, id int, reason string, blockerId *int, blockedUntil *time.Time) error
	IsUserBlocked(ctx context.Context, id int) (bool, error)
	ModifyPassword(ctx context.Context, id int, password string) error
	AddNotificationToken(ctx context.Context, id int, text string) error
	GetUserNotificationsToken(ctx context.Context, id int) (models.Notifications, error)
	SendNotifByMobile(cont context.Context, userId int, notification models.NotifyRequest) error
	SendNotifByEmail(cont context.Context, userId int, request models.NotifyRequest) error
	VerifyUser(ctx context.Context, id int) error
}

type userService struct {
	userRepo      repo.UserRepository
	blockUserRepo repo.BlockedUserRepository
}

func NewUserService(userRepo repo.UserRepository, blockedUserRepo repo.BlockedUserRepository) *userService {
	return &userService{userRepo: userRepo, blockUserRepo: blockedUserRepo}
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, request models.CreateUserRequest) (int, error) {

	hashPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return 0, err
	}

	user := &models.User{
		Email:    request.Email,
		Password: hashPassword,
		Name:     request.Name,
		Surname:  request.Surname,
		Role:     request.Role,
		Verified: request.Verified,
	}

	return s.userRepo.AddUser(ctx, user)
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

func (s *userService) ModifyUser(ctx context.Context, id int, user *models.User) error {
	tableUser, err := s.userRepo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	// Update the existing user with the new values
	tableUser.Update(user)

	return s.userRepo.ModifyUser(ctx, tableUser)
}

func (s *userService) IsUserBlocked(ctx context.Context, id int) (bool, error) {
	blocks, err := s.blockUserRepo.GetBlocksByUserId(ctx, id)
	if err != nil {
		return false, err
	}

	return len(blocks) > 0, nil
}

func (s *userService) BlockUser(
	ctx context.Context,
	userId int,
	reason string,
	blockerId *int,
	blockedUntil *time.Time,
) error {
	if _, err := s.userRepo.GetUser(ctx, userId); err != nil {
		return err
	}

	if err := s.blockUserRepo.BlockUser(ctx, userId, reason, blockerId, blockedUntil); err != nil {
		return err
	}

	return nil
}
func (s *userService) ModifyPassword(ctx context.Context, id int, password string) error {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return s.userRepo.ModifyPassword(ctx, id, hashPassword)
}

func (s *userService) AddNotificationToken(ctx context.Context, id int, text string) error {
	return s.userRepo.AddNotificationToken(ctx, id, text)
}
func (s *userService) GetUserNotificationsToken(ctx context.Context, id int) (models.Notifications, error) {
	return s.userRepo.GetUserNotificationsToken(ctx, id)
}

func (s *userService) VerifyUser(ctx context.Context, id int) error {
	return s.userRepo.SetVerifiedTrue(ctx, id)
}

func (s *userService) SendNotifByMobile(cont context.Context, userId int, notification models.NotifyRequest) error {
	return nil
}

func (s *userService) SendNotifByEmail(cont context.Context, userId int, request models.NotifyRequest) error {
	user, err := s.userRepo.GetUser(cont, userId)
	if err != nil {
		return err
	}
	from := mail.NewEmail("ClassConnect service", "bmorseletto@fi.uba.ar")
	subject := request.NotificationTitle
	to := mail.NewEmail("User", user.Email)
	content := mail.NewContent("text/plain", request.NotificationText)
	message := mail.NewV3MailInit(from, subject, to, content)

	client := sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY"))
	_, err2 := client.Send(message)
	return err2
}
