package services_test

import (
	"context"
	"errors"
	"testing"

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
	service := services.NewUserService(mockRepo, mockBlockedRepo)

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
	service := services.NewUserService(mockRepo, mockBlockedRepo)

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
	service := services.NewUserService(mockRepo, mockBlockedRepo)

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
	service := services.NewUserService(mockRepo, mockBlockedRepo)

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
	user, err := service.CreateUser(ctx, createRequest, false)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, user)
}

func TestUserService_DeleteUser(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	mockRepo.EXPECT().DeleteUser(ctx, 1).Return(nil)

	// Act
	err := service.DeleteUser(ctx, 1)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_ModifyUser(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo)

	userToModify := &models.User{
		Id:       1,
		Name:     "John Updated",
		Surname:  "Doe Updated",
		Email:    "john@example.com",
		Location: "New York",
	}

	ctx := context.Background()
	mockRepo.EXPECT().ModifyUser(ctx, userToModify).Return(nil)

	// Act
	err := service.ModifyUser(ctx, userToModify)

	// Assert
	assert.NoError(t, err)
}

func TestUserService_ModifyLocation(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo)

	ctx := context.Background()
	mockRepo.EXPECT().ModifyLocation(ctx, 1, "New York").Return(nil)

	// Act
	err := service.ModifyLocation(ctx, 1, "New York")

	// Assert
	assert.NoError(t, err)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	// Arrange
	mockRepo := repositories.NewMockUserRepository(t)
	mockBlockedRepo := repositories.NewMockBlockedUserRepository(t)
	service := services.NewUserService(mockRepo, mockBlockedRepo)

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

	service := services.NewUserService(mockUserRepo, mockBlockedRepo)

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
	service := services.NewUserService(mockUserRepo, mockBlockedRepo)

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
	mockRepo := NewMockUserRepository(t)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("ModifyPassword", ctx, 1, mock.Anything).Return(nil)

	err := service.ModifyPassword(ctx, 1, TEST_PASSWORD)
	assert.NoError(t, err)
}
