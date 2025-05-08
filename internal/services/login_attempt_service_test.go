package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginAttemptService_AddLoginAttempt_Successful(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	successful := true

	mockRepo.EXPECT().AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful).Return(nil)

	// Act
	err := service.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)

	// Assert
	assert.NoError(t, err)
}

func TestLoginAttemptService_AddLoginAttempt_Failed_NoBlock(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	successful := false

	mockRepo.EXPECT().AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful).Return(nil)
	mockRepo.EXPECT().BadLoginAttemptsInLast10Minutes(ctx, userID).Return(services.MaxFailedAttempts-1, nil)

	// Act
	err := service.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)

	// Assert
	assert.NoError(t, err)
}

func TestLoginAttemptService_AddLoginAttempt_Failed_WithBlock(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	successful := false

	mockRepo.EXPECT().AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful).Return(nil)
	mockRepo.EXPECT().BadLoginAttemptsInLast10Minutes(ctx, userID).Return(services.MaxFailedAttempts, nil)

	// We use nil literal instead of *int or *time.Time to match the implementation's parameter types
	mockBlockedRepo.EXPECT().BlockUser(
		ctx,
		userID,
		services.BlockedUserReason,
		mock.Anything,
		mock.Anything,
	).Return(nil)

	// Act
	err := service.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)

	// Assert
	assert.NoError(t, err)
}

func TestLoginAttemptService_AddLoginAttempt_Error(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	successful := true

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful).Return(expectedErr)

	// Act
	err := service.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLoginAttemptService_GetLoginsByUserId(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	limit := 10
	offset := 0

	expectedLogins := []*models.LoginAttempt{
		{
			ID:         1,
			UserID:     userID,
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			Successful: true,
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         2,
			UserID:     userID,
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			Successful: false,
			CreatedAt:  time.Now(),
		},
	}

	mockRepo.EXPECT().GetLoginsByUserId(ctx, userID, limit, offset).Return(expectedLogins, nil)

	// Act
	logins, err := service.GetLoginsByUserId(ctx, userID, limit, offset)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLogins, logins)
}

func TestLoginAttemptService_GetLoginsByUserId_Error(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockLoginAttemptRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewLoginAttemptService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	userID := 1
	limit := 10
	offset := 0

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().GetLoginsByUserId(ctx, userID, limit, offset).Return(nil, expectedErr)

	// Act
	logins, err := service.GetLoginsByUserId(ctx, userID, limit, offset)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, logins)
}
