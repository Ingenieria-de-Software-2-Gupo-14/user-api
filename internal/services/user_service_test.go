package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/utils"
	"ing-soft-2-tp1/mocks"
	"testing"
)

const (
	TEST_ID              = 1
	TEST_USERNAME        = "testUser"
	TEST_NAME            = "testName"
	TEST_SURNAME         = "testSurname"
	TEST_PASSWORD        = "testPassword"
	TEST_EMAIL           = "testEmail"
	TEST_LOCATION        = "testLocation"
	TEST_ADMIN           = false
	TEST_BLOCKED         = false
	TEST_PROFILE_PICTURE = 0
	TEST_DESCRIPTION     = "testDesc"
)

var ExpectedUser = models.User{
	Username:     TEST_USERNAME,
	Name:         TEST_NAME,
	Surname:      TEST_SURNAME,
	Password:     TEST_PASSWORD,
	Email:        TEST_EMAIL,
	Location:     TEST_LOCATION,
	Admin:        TEST_ADMIN,
	BlockedUser:  TEST_BLOCKED,
	ProfilePhoto: TEST_PROFILE_PICTURE,
	Description:  TEST_DESCRIPTION,
}

func TestNewUserService(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	result := NewUserService(mockRepo)

	assert.NotNil(t, result)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("DeleteUser", ctx, TEST_ID).Return(nil)

	err := service.DeleteUser(ctx, TEST_ID)
	assert.NoError(t, err)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}

	expectedUser := models.User{
		Id:       1,
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}

	mockRepo.On("AddUser", ctx, mock.Anything).Return(1, nil)

	user, err := service.CreateUser(ctx, request, false)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Surname, user.Surname)
	hashErr := utils.CompareHashPassword(user.Password, expectedUser.Password)
	assert.NoError(t, hashErr)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Location, user.Location)
	assert.Equal(t, expectedUser.Admin, user.Admin)
	assert.Equal(t, expectedUser.BlockedUser, user.BlockedUser)
	assert.Equal(t, expectedUser.ProfilePhoto, user.ProfilePhoto)
	assert.Equal(t, expectedUser.Description, user.Description)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetUserByEmail", ctx, TEST_EMAIL).Return(&ExpectedUser, nil)

	user, err := service.GetUserByEmail(ctx, TEST_EMAIL)
	assert.NoError(t, err)
	assert.Equal(t, ExpectedUser, *user)
}

func TestUserService_GetUserById(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	result := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("GetUser", ctx, TEST_ID).Return(&ExpectedUser, nil)

	user, err := result.GetUserById(ctx, TEST_ID)
	assert.NoError(t, err)
	assert.Equal(t, ExpectedUser, *user)
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	expectedUsers := []models.User{}
	expectedUsers = append(expectedUsers, ExpectedUser)

	mockRepo.On("GetAllUsers", ctx).Return(expectedUsers, nil)

	users, err := service.GetAllUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestUserService_BlockUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("BlockUser", ctx, 1).Return(nil)

	err := service.BlockUser(ctx, 1)
	assert.NoError(t, err)
}

func TestUserService_ModifyUser(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("ModifyUser", ctx, &ExpectedUser).Return(nil)

	err := service.ModifyUser(ctx, &ExpectedUser)
	assert.NoError(t, err)
}

func TestUserService_ModifyLocation(t *testing.T) {
	mockRepo := new(mocks.UserRepository)

	service := NewUserService(mockRepo)

	ctx := context.Background()

	mockRepo.On("ModifyLocation", ctx, 1, TEST_LOCATION).Return(nil)

	err := service.ModifyLocation(ctx, 1, TEST_LOCATION)
	assert.NoError(t, err)
}
