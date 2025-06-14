package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/controller"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	s "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (*s.MockUserService, *gin.Context, *httptest.ResponseRecorder, *controller.UserController) {
	gin.SetMode(gin.TestMode)
	mockService := s.NewMockUserService(t)
	mockRulesService := s.NewMockRulesService(t)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	userController := controller.CreateController(mockService, mockRulesService)
	return mockService, c, recorder, userController
}

func TestCreateController(t *testing.T) {
	mockService := s.NewMockUserService(t)
	mockRulesService := s.NewMockRulesService(t)
	result := controller.CreateController(mockService, mockRulesService)
	assert.NotNil(t, result)
}

func TestUsersGet_Success(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	c.Request = httptest.NewRequest(http.MethodGet, "/api/users", nil)

	// Setup expected users
	expectedUsers := []models.User{
		{
			Id:      1,
			Name:    "Test1",
			Surname: "User1",
			Email:   "test1@example.com",
		},
		{
			Id:      2,
			Name:    "Test2",
			Surname: "User2",
			Email:   "test2@example.com",
		},
	}

	mockService.EXPECT().GetAllUsers(mock.Anything).Return(expectedUsers, nil)

	// Call the function
	userController.UsersGet(c)

	// Check response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response struct {
		Data []models.User `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedUsers), len(response.Data))
	assert.Equal(t, expectedUsers[0].Id, response.Data[0].Id)
	assert.Equal(t, expectedUsers[0].Name, response.Data[0].Name)
}

func TestUserGetById_Success(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	userId := 1
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users/"+strconv.Itoa(userId), nil)
	c.AddParam("id", strconv.Itoa(userId))

	expectedUser := &models.User{
		Id:      userId,
		Name:    "Test",
		Surname: "User",
		Email:   "test@example.com",
	}

	mockService.EXPECT().GetUserById(mock.Anything, userId).Return(expectedUser, nil)

	// Call the function
	userController.UserGetById(c)

	// Check response
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response struct {
		Data models.User `json:"data"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, response.Data.Id)
	assert.Equal(t, expectedUser.Name, response.Data.Name)
}

func TestUserGetById_NotFound(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	userId := 999
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users/"+strconv.Itoa(userId), nil)
	c.AddParam("id", strconv.Itoa(userId))

	mockService.EXPECT().GetUserById(mock.Anything, userId).Return(nil, errors.New("not found"))

	// Call the function
	userController.UserGetById(c)

	// Check response
	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUserDeleteById_Success(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	userId := 1
	c.Request = httptest.NewRequest(http.MethodDelete, "/api/users/"+strconv.Itoa(userId), nil)
	c.AddParam("id", strconv.Itoa(userId))

	mockService.EXPECT().DeleteUser(mock.Anything, userId).Return(nil)

	// Call the function
	userController.UserDeleteById(c)

	// Check response
	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestModifyUser_Success(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	updatedUserDto := models.UserUpdateDto{
		Name:     "Updated",
		Surname:  "User",
		Location: "New Location",
	}

	jsonValue, _ := json.Marshal(updatedUserDto)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/users/1", bytes.NewBuffer(jsonValue))
	c.Request.Header.Set("Content-Type", "application/json")
	c.AddParam("id", "1")

	mockService.EXPECT().ModifyUser(mock.Anything, 1, mock.MatchedBy(func(u models.UserUpdateDto) bool {
		return u.Name == updatedUserDto.Name &&
			u.Surname == updatedUserDto.Surname &&
			u.Location == updatedUserDto.Location
	})).Return(nil)

	// Call the function
	userController.ModifyUser(c)

	// Check response
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestBlockUserById_Success(t *testing.T) {
	mockService, c, recorder, userController := setupTest(t)

	userId := 1
	c.Request = httptest.NewRequest(http.MethodPost, "/api/users/"+strconv.Itoa(userId)+"/block", nil)
	c.AddParam("id", strconv.Itoa(userId))
	c.Request = c.Request.WithContext(context.Background())

	mockService.EXPECT().BlockUser(mock.Anything, userId, "", mock.AnythingOfType("*int"), mock.AnythingOfType("*time.Time")).Return(nil)

	// Call the function
	userController.BlockUserById(c)

	// Check response
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUserController_ModifyPassword(t *testing.T) {
	mockService, c, recorder, controller := setupTest(t)

	newPassword := "TEST_PASSWORD"

	gin.SetMode(gin.TestMode)

	expectedRequest := models.PasswordModifyRequest{
		Token:    "123456",
		Password: newPassword,
	}
	expectedPasswordResetData := models.PasswordResetData{
		Email:  "TEST_EMAIL",
		UserId: 1,
		Exp:    time.Now(),
		Used:   false,
	}

	jsonBody, _ := json.Marshal(expectedRequest)

	req, _ := http.NewRequest(http.MethodPut, "/users/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockService.On("ValidatePasswordResetToken", c.Request.Context(), expectedRequest.Token).Return(&expectedPasswordResetData, nil)
	mockService.On("ModifyPassword", c.Request.Context(), 1, newPassword).Return(nil)
	mockService.On("SetPasswordTokenUsed", c.Request.Context(), expectedRequest.Token).Return(nil)

	controller.ModifyUserPasssword(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestUserController_ModifyPassword_WrongParam(t *testing.T) {
	_, c, recorder, controller := setupTest(t)

	newPassword := "TEST_PASSWORD"

	gin.SetMode(gin.TestMode)

	expectedRequest := models.PasswordModifyRequest{
		Password: newPassword,
	}

	jsonBody, _ := json.Marshal(expectedRequest)

	req, _ := http.NewRequest(http.MethodPost, "/users/a/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	c.Request = req

	controller.ModifyUserPasssword(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	var result models.ErrorResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
