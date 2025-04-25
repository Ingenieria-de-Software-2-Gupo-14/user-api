package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ing-soft-2-tp1/internal/errors"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/utils"
	"ing-soft-2-tp1/mocks"
	"net/http"
	"net/http/httptest"
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

func TestCreateController(t *testing.T) {
	mockService := new(mocks.UserService)

	result := CreateController(mockService)

	assert.NotNil(t, result)
}

func TestUserController_Health(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, "/health", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Health(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_RegisterUser(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(nil, errors.ErrNotFound)
	mockService.On("CreateUser", c.Request.Context(), request, false).Return(&expectedUser, nil)

	controller.RegisterUser(c)

	var result models.ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ResponseUser{User: expectedUser}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_RegisterUser_UserAlreadyExists(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(&expectedUser, nil)

	controller.RegisterUser(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleConflict,
		Status:   http.StatusConflict,
		Detail:   models.ErrorDescriptionConflict,
		Instance: "/users",
	}

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_RegisterUser_InternalError(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(nil, errors.ErrNotFound)
	mockService.On("CreateUser", c.Request.Context(), request, false).Return(nil, sql.ErrConnDone)

	controller.RegisterUser(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UsersGet(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	expectedUsers := []models.User{}
	expectedUsers = append(expectedUsers, expectedUser)

	mockService.On("GetAllUsers", c.Request.Context()).Return(expectedUsers, nil)

	controller.UsersGet(c)

	var result models.ResponseUsers
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResult := models.ResponseUsers{Users: expectedUsers}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResult, result)
}

func TestUserController_UserGetById_NotFound(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("GetUserById", c.Request.Context(), 1).Return(nil, errors.ErrNotFound)

	controller.UserGetById(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleUserNotFound,
		Status:   http.StatusNotFound,
		Detail:   models.ErrorDescriptionUserNotFound,
		Instance: "/users/1",
	}

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UserGetById(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserById", c.Request.Context(), 1).Return(&expectedUser, nil)

	controller.UserGetById(c)

	var result models.ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ResponseUser{User: expectedUser}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UserDeleteById(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("DeleteUser", c.Request.Context(), 1).Return(nil)

	controller.UserDeleteById(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserController_UserDeleteById_InternalError(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("DeleteUser", c.Request.Context(), 1).Return(sql.ErrConnDone)

	controller.UserDeleteById(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/1",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UserDeleteById_WrongParams(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodDelete, "/users/a", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	c.Request = req

	mockService.On("DeleteUser", c.Request.Context(), 1).Return(sql.ErrConnDone)

	controller.UserDeleteById(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/a",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_RegisterAdmin(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/admins", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(nil, errors.ErrNotFound)
	mockService.On("CreateUser", c.Request.Context(), request, true).Return(&expectedUser, nil)

	controller.RegisterAdmin(c)

	var result models.ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ResponseUser{User: expectedUser}

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedResponse, result)

}

func TestUserController_RegisterAdmin_UserAlreadyExists(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.CreateUserRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
		Name:     TEST_NAME,
		Surname:  TEST_SURNAME,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/admins", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        TEST_ADMIN,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(&expectedUser, nil)

	controller.RegisterAdmin(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleConflict,
		Status:   http.StatusConflict,
		Detail:   models.ErrorDescriptionConflict,
		Instance: "/admins",
	}

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UserLogin(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.LoginRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	hashPassword, _ := utils.HashPassword(TEST_PASSWORD)

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     hashPassword,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(&expectedUser, nil)

	controller.UserLogin(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_UserLogin_WrongPassword(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.LoginRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     TEST_PASSWORD,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(&expectedUser, nil)

	controller.UserLogin(c)
	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    "error",
		Status:   http.StatusUnauthorized,
		Detail:   "error",
		Instance: "/login",
	}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_UserLogin_NoUser(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.LoginRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(nil, errors.ErrNotFound)

	controller.UserLogin(c)
	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleUserNotFound,
		Status:   http.StatusNotFound,
		Detail:   models.ErrorDescriptionUserNotFound,
		Instance: "/login",
	}

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_ModifyUser(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     "",
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	jsonBody, _ := json.Marshal(expectedUser)

	req, _ := http.NewRequest(http.MethodPost, "/users/modify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("ModifyUser", c.Request.Context(), &expectedUser).Return(nil)

	controller.ModifyUser(c)

	var result models.ResponseUser
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ResponseUser{User: expectedUser}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_ModifyUser_InternalError(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     "",
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	jsonBody, _ := json.Marshal(expectedUser)

	req, _ := http.NewRequest(http.MethodPost, "/users/modify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockService.On("ModifyUser", c.Request.Context(), &expectedUser).Return(sql.ErrConnDone)

	controller.ModifyUser(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/modify",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_ModifyUserLocation(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	expectedUser := models.User{
		Location: TEST_LOCATION,
	}
	jsonBody, _ := json.Marshal(expectedUser)
	req, _ := http.NewRequest(http.MethodPut, "/users/1/location", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("ModifyLocation", c, 1, TEST_LOCATION).Return(nil)

	controller.ModifyUserLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_ModifyUserLocation_InternalError(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	expectedUser := models.User{
		Location: TEST_LOCATION,
	}
	jsonBody, _ := json.Marshal(expectedUser)
	req, _ := http.NewRequest(http.MethodPut, "/users/1/location", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("ModifyLocation", c, 1, TEST_LOCATION).Return(sql.ErrConnDone)

	controller.ModifyUserLocation(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/1/location",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_ModifyUserLocation_WrongParams(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	expectedUser := models.User{
		Location: TEST_LOCATION,
	}
	jsonBody, _ := json.Marshal(expectedUser)
	req, _ := http.NewRequest(http.MethodPut, "/users/a/location", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	c.Request = req

	controller.ModifyUserLocation(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/a/location",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_BlockUserById(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPut, "/users/block/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("BlockUser", c, 1).Return(nil)

	controller.BlockUserById(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_BlockUserById_InternalError(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPut, "/users/block/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockService.On("BlockUser", c, 1).Return(sql.ErrConnDone)

	controller.BlockUserById(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/block/1",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_BlockUserById_WrongParam(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPut, "/users/block/a", nil)
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	c.Request = req

	controller.BlockUserById(c)

	var result models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		return
	}

	expectedResponse := models.ErrorResponse{
		Type:     models.ErrorTypeBlank,
		Title:    models.ErrorTitleInternalServerError,
		Status:   http.StatusInternalServerError,
		Detail:   models.ErrorDescriptionInternalServerError,
		Instance: "/users/block/a",
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponse, result)
}

func TestUserController_ValidateToken(t *testing.T) {
	mockService := new(mocks.UserService)

	controller := CreateController(mockService)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	request := models.LoginRequest{
		Email:    TEST_EMAIL,
		Password: TEST_PASSWORD,
	}
	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	hashPassword, _ := utils.HashPassword(TEST_PASSWORD)

	expectedUser := models.User{
		Id:           1,
		Username:     "",
		Name:         TEST_NAME,
		Surname:      TEST_SURNAME,
		Password:     hashPassword,
		Email:        TEST_EMAIL,
		Location:     "",
		Admin:        true,
		BlockedUser:  TEST_BLOCKED,
		ProfilePhoto: TEST_PROFILE_PICTURE,
		Description:  "",
	}

	mockService.On("GetUserByEmail", c.Request.Context(), TEST_EMAIL).Return(&expectedUser, nil)

	controller.UserLogin(c)
	controller.ValidateToken(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
