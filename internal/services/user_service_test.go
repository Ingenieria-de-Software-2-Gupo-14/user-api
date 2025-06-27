package services_test

import (
	"context"
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/sendgrid/rest"
	"testing"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetAllUsers(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	expectedUsers := []models.User{
		{Id: 1, Name: "John", Surname: "Doe", Email: "john@example.com"},
		{Id: 2, Name: "Jane", Surname: "Doe", Email: "jane@example.com"},
	}

	ctx := context.Background()
	mockRepo.EXPECT().GetAllUsers(ctx).Return(expectedUsers, nil)

	// Act
	users, err := service.GetAllUsers(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestUserService_GetUserById(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	expectedUser := &models.User{
		Id:      1,
		Name:    "John",
		Surname: "Doe",
		Email:   "john@example.com",
	}

	ctx := context.Background()
	mockRepo.EXPECT().GetUser(ctx, 1).Return(expectedUser, nil)

	// Act
	user, err := service.GetUserById(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_GetUserById_Error(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	expectedErr := errors.New("user not found")
	ctx := context.Background()
	mockRepo.EXPECT().GetUser(ctx, 999).Return(nil, expectedErr)

	// Act
	user, err := service.GetUserById(ctx, 999)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, user)
}

func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	createRequest := models.CreateUserRequest{
		Name:     "John",
		Surname:  "Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	ctx := context.Background()

	// Mock AddUser to return user ID
	mockRepo.EXPECT().AddUser(ctx, mock.AnythingOfType("*models.User")).Return(1, nil)

	// Act
	user, err := service.CreateUser(ctx, createRequest)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, user)
}

func TestUserService_DeleteUser(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	mockRepo.EXPECT().DeleteUser(ctx, 1).Return(nil)

	// Act
	err := service.DeleteUser(ctx, 1)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_MakeTeacher(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	mockRepo.EXPECT().MakeTeacher(ctx, 1).Return(nil)

	// Act
	err := service.MakeTeacher(ctx, 1)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_ModifyUser(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	userId := 1
	userToModify := models.UserUpdateDto{
		Name:     "John Updated",
		Surname:  "Doe Updated",
		Location: "New York",
	}

	existingUser := &models.User{
		Id:       userId,
		Name:     "John",
		Surname:  "Doe",
		Location: "Old Location",
	}

	ctx := context.Background()
	mockRepo.EXPECT().GetUser(ctx, userId).Return(existingUser, nil)
	mockRepo.EXPECT().ModifyUser(ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.Id == userId &&
			u.Name == "John Updated" &&
			u.Surname == "Doe Updated" &&
			u.Location == "New York"
	})).Return(nil)

	// Act
	err := service.ModifyUser(ctx, userId, userToModify)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_ModifyLocation(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userId := 1
	newLocation := "New York"

	existingUser := &models.User{
		Id:       userId,
		Name:     "John",
		Surname:  "Doe",
		Email:    "john@example.com",
		Location: "Old Location",
	}

	// Mock GetUser to return the existing user
	mockRepo.EXPECT().GetUser(ctx, userId).Return(existingUser, nil)

	// Create a copy of the existing user and update it with the new location
	expectedUpdatedUser := *existingUser
	expectedUpdatedUser.Location = newLocation

	// Mock ModifyUser with the expected updated user
	mockRepo.EXPECT().ModifyUser(ctx, existingUser).Return(nil)

	// Create a user with only the location field set for the input
	locationUser := models.UserUpdateDto{
		Location: newLocation,
	}

	// Act
	err := service.ModifyUser(ctx, userId, locationUser)

	// Assert
	assert.NoError(t, err)
	// Verify that the existing user was updated with the new location
	assert.Equal(t, newLocation, existingUser.Location)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	expectedUser := &models.User{
		Id:      1,
		Name:    "John",
		Surname: "Doe",
		Email:   "john@example.com",
	}

	ctx := context.Background()
	mockRepo.EXPECT().GetUserByEmail(ctx, "john@example.com").Return(expectedUser, nil)

	// Act
	user, err := service.GetUserByEmail(ctx, "john@example.com")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestUserService_BlockUser(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockUserRepo, mockBlockedRepo, mockEmail)

	expectedUser := &models.User{
		Id:      1,
		Name:    "John",
		Surname: "Doe",
		Email:   "john@example.com",
	}

	ctx := context.Background()
	mockUserRepo.EXPECT().GetUser(ctx, 1).Return(expectedUser, nil)
	mockBlockedRepo.EXPECT().BlockUser(ctx, 1, "", mock.Anything, mock.Anything).Return(nil)

	// Act
	err := service.BlockUser(ctx, 1, "", nil, nil)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_IsUserBlocked(t *testing.T) {
	// Arrange
	mockUserRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockUserRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userId := 1

	mockBlockedRepo.EXPECT().GetBlocksByUserId(ctx, userId).Return([]models.BlockedUser{}, nil)

	// Act
	blocked, err := service.IsUserBlocked(ctx, userId)

	// Assert
	assert.NoError(t, err)
	assert.False(t, blocked)
}

func TestUserService_ModifyPassword(t *testing.T) {
	newPassword := "TEST_PASSWORD"

	mockUserRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockUserRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()

	mockUserRepo.On("ModifyPassword", ctx, 1, mock.Anything).Return(nil)

	err := service.ModifyPassword(ctx, 1, newPassword)
	assert.NoError(t, err)
}

func TestUserService_AddNotificationToken(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	token := "123456"

	mockRepo.EXPECT().
		AddNotificationToken(ctx, userID, token).
		Return(nil)

	err := service.AddNotificationToken(ctx, userID, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserNotificationsToken(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	notificationToken := models.NotificationToken{
		NotificationToken: "123456",
		CreatedTime:       time.Time{},
	}
	expectedTokens := models.NotificationTokens{
		NotificationTokens: []models.NotificationToken{notificationToken},
	}

	mockRepo.EXPECT().
		GetUserNotificationsToken(ctx, userID).
		Return(expectedTokens, nil)

	tokens, err := service.GetUserNotificationsToken(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTokens, tokens)
	mockRepo.AssertExpectations(t)
}

func TestUserService_VerifyUser(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1

	mockRepo.EXPECT().
		SetVerifiedTrue(ctx, userID).
		Return(nil)

	err := service.VerifyUser(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_SetPasswordTokenUsed(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	token := "123456"

	mockRepo.EXPECT().
		SetPasswordTokenUsed(ctx, token).
		Return(nil)

	err := service.SetPasswordTokenUsed(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_SetNotificationPreference(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	preference := models.NotificationPreferenceRequest{
		NotificationType:       "exam_notification",
		NotificationPreference: false,
	}

	mockRepo.EXPECT().
		SetNotificationPreference(ctx, userID, preference).
		Return(nil)

	err := service.SetNotificationPreference(ctx, userID, preference)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CheckPreference(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	notificationType := "email"

	mockRepo.EXPECT().
		CheckPreference(ctx, userID, notificationType).
		Return(true, nil)

	result, err := service.CheckPreference(ctx, userID, notificationType)

	assert.NoError(t, err)
	assert.True(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetNotificationPreference(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	expected := &models.NotificationPreference{
		ExamNotification:     true,
		HomeworkNotification: false,
		SocialNotification:   false,
	}

	mockRepo.EXPECT().
		GetNotificationPreference(ctx, userID).
		Return(expected, nil)

	result, err := service.GetNotificationPreference(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ValidatePasswordResetToken_Valid(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	token := "valid-token"
	expected := &models.PasswordResetData{
		Email:  "test@email.com",
		UserId: 1,
		Exp:    time.Now().Add(1 * time.Hour),
		Used:   false,
	}

	mockRepo.EXPECT().
		GetPasswordResetTokenInfo(ctx, token).
		Return(expected, nil)

	result, err := service.ValidatePasswordResetToken(ctx, token)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ValidatePasswordResetToken_Expired(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	token := "expired-token"
	expired := &models.PasswordResetData{
		Email:  "test@email.com",
		UserId: 1,
		Exp:    time.Now().Add(-1 * time.Hour),
		Used:   false,
	}

	mockRepo.EXPECT().
		GetPasswordResetTokenInfo(ctx, token).
		Return(expired, nil)

	result, err := service.ValidatePasswordResetToken(ctx, token)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "Token Expired")
	mockRepo.AssertExpectations(t)
}

func TestSendNotifByEmail(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userId := 1
	email := "user@test.com"

	user := &models.User{
		Id:    userId,
		Email: email,
	}

	mockRepo.EXPECT().
		GetUser(ctx, userId).
		Return(user, nil)

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	req := models.NotifyRequest{
		NotificationTitle: "test",
		NotificationText:  "test message",
	}

	err := service.SendNotifByEmail(ctx, userId, req)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestStartPasswordReset(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	email := "user@example.com"

	mockRepo.EXPECT().
		AddPasswordResetToken(ctx, userID, email, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
		Return(nil)

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	err := service.StartPasswordReset(ctx, userID, email)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestSendNotifByMobile_Expo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://exp.host/--/api/v2/push/send",
		httpmock.NewStringResponder(200, `{"status":"ok"}`),
	)

	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	mockRepo.
		EXPECT().
		GetUserNotificationsToken(context.Background(), 1).
		Return(models.NotificationTokens{
			NotificationTokens: []models.NotificationToken{
				{NotificationToken: "ExponentPushToken[1234567890]"},
			},
		}, nil)

	err := service.SendNotifByMobile(context.Background(), 1, models.NotifyRequest{
		NotificationTitle: "Test Title",
		NotificationText:  "This is a test",
	})

	assert.NoError(t, err)

}

func TestSendNotifByMobile_FirebaseServiceAccountNotSet(t *testing.T) {
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	mockEmail := services.NewMockEmailSender(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo, mockEmail)

	t.Setenv("FIREBASE_SERVICE_ACCOUNT", "")

	userId := 1
	notification := models.NotifyRequest{
		NotificationTitle: "Test Title",
		NotificationText:  "Test Body",
	}

	// Non-Expo token to trigger Firebase path
	mockRepo.EXPECT().
		GetUserNotificationsToken(context.Background(), userId).
		Return(models.NotificationTokens{
			NotificationTokens: []models.NotificationToken{
				{NotificationToken: "firebase_token_123"},
			},
		}, nil).
		Once()

	err := service.SendNotifByMobile(context.Background(), userId, notification)

	assert.NoError(t, err)
}
