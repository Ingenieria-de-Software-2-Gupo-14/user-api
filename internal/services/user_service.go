package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/oauth2/google"
	"net/http"
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
	ModifyUser(ctx context.Context, id int, user models.UserUpdateDto) error
	BlockUser(ctx context.Context, id int, reason string, blockerId *int, blockedUntil *time.Time) error
	IsUserBlocked(ctx context.Context, id int) (bool, error)
	ModifyPassword(ctx context.Context, id int, password string) error
	AddNotificationToken(ctx context.Context, id int, text string) error
	GetUserNotificationsToken(ctx context.Context, id int) (models.NotificationTokens, error)
	SendNotifByMobile(cont context.Context, userId int, notification models.NotifyRequest) error
	SendNotifByEmail(cont context.Context, userId int, request models.NotifyRequest) error
	VerifyUser(ctx context.Context, id int) error
	StartPasswordReset(ctx context.Context, id int, email string) error
	ValidatePasswordResetToken(ctx context.Context, token string) (*models.PasswordResetData, error)
	SetPasswordTokenUsed(ctx context.Context, token string) error
	SetNotificationPreference(ctx context.Context, id int, preference models.NotificationPreferenceRequest) error
	CheckPreference(ctx context.Context, id int, notificationType string) (bool, error)
	GetNotificationPreference(ctx context.Context, id int) (*models.NotificationPreference, error)
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

func (s *userService) ModifyUser(ctx context.Context, id int, user models.UserUpdateDto) error {
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
func (s *userService) GetUserNotificationsToken(ctx context.Context, id int) (models.NotificationTokens, error) {
	return s.userRepo.GetUserNotificationsToken(ctx, id)
}

func (s *userService) VerifyUser(ctx context.Context, id int) error {
	return s.userRepo.SetVerifiedTrue(ctx, id)
}

func (s *userService) SendNotifByMobile(cont context.Context, userId int, notification models.NotifyRequest) error {
	tokens, _ := s.GetUserNotificationsToken(cont, userId)
	for _, token := range tokens.NotificationTokens {
		err := sendNotifToDevice(token.NotificationToken, notification)
		if err != nil {
			println(err.Error())
		}
	}
	return nil
}

func sendNotifToDevice(userToken string, notification models.NotifyRequest) error {
	svcJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	if svcJSON == "" {
		return fmt.Errorf("FIREBASE_SERVICE_ACCOUNT not set")
	}
	creds := []byte(svcJSON)

	conf, err := google.JWTConfigFromJSON(creds, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return fmt.Errorf("invalid service account JSON: %v", err)
	}

	client := conf.Client(context.Background())

	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", os.Getenv("FCM_PROJECT_ID"))
	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"token": userToken,
			"notification": map[string]string{
				"title": notification.NotificationTitle,
				"body":  notification.NotificationText,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending FCM request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response from FCM: %v", resp.Status)
	}
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

func (s *userService) StartPasswordReset(ctx context.Context, id int, email string) error {
	token, err := password.Generate(6, 2, 0, false, true)
	if err != nil {
		return err
	}
	err = s.userRepo.AddPasswordResetToken(ctx, id, email, token, time.Now().Add(5*time.Minute))
	if err != nil {
		return err
	}
	from := mail.NewEmail("ClassConnect service", "bmorseletto@fi.uba.ar")
	subject := "Password Reset"
	to := mail.NewEmail("User", email)
	resetLink := "https://user-api-production-99c2.up.railway.app/users/reset/password?token=" + token
	plainTextContent := resetLink
	htmlContent := fmt.Sprintf(`
	<html>
		<body>
			<p>Hello,</p>
			<p>Click the link below to reset your password:</p>
			<p><a href="%s">Reset password</a></p>
		</body>
	</html>`, resetLink)

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject

	p := mail.NewPersonalization()
	p.AddTos(to)
	message.AddPersonalizations(p)

	message.AddContent(mail.NewContent("text/plain", plainTextContent))
	message.AddContent(mail.NewContent("text/html", htmlContent))

	client := sendgrid.NewSendClient(os.Getenv("EMAIL_API_KEY"))
	_, err2 := client.Send(message)
	return err2
}

func (s *userService) ValidatePasswordResetToken(ctx context.Context, token string) (*models.PasswordResetData, error) {
	info, err := s.userRepo.GetPasswordResetTokenInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	if info.Exp.Before(time.Now()) {
		return nil, errors.New("Token Expired")
	}
	return info, nil
}

func (s *userService) SetPasswordTokenUsed(ctx context.Context, token string) error {
	return s.userRepo.SetPasswordTokenUsed(ctx, token)
}

func (s *userService) SetNotificationPreference(ctx context.Context, id int, preference models.NotificationPreferenceRequest) error {
	return s.userRepo.SetNotificationPreference(ctx, id, preference)
}

func (s *userService) CheckPreference(ctx context.Context, id int, notificationType string) (bool, error) {
	return s.userRepo.CheckPreference(ctx, id, notificationType)
}

func (s *userService) GetNotificationPreference(ctx context.Context, id int) (*models.NotificationPreference, error) {
	return s.userRepo.GetNotificationPreference(ctx, id)
}
