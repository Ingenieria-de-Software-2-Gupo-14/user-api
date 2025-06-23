package controller_test

import (
	"bytes"
	"context"
	"database/sql"
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

func setupTest(t *testing.T) (*s.MockUserService, *s.MockRulesService, *gin.Context, *httptest.ResponseRecorder, *controller.UserController) {
	gin.SetMode(gin.TestMode)
	mockService := s.NewMockUserService(t)
	mockRulesService := s.NewMockRulesService(t)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	userController := controller.CreateController(mockService, mockRulesService)
	return mockService, mockRulesService, c, recorder, userController
}

func TestCreateController(t *testing.T) {
	mockService := s.NewMockUserService(t)
	mockRulesService := s.NewMockRulesService(t)
	result := controller.CreateController(mockService, mockRulesService)
	assert.NotNil(t, result)
}

func TestUsersGet_Success(t *testing.T) {
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, userController := setupTest(t)

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
	mockService, _, c, recorder, controller := setupTest(t)

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
	_, _, c, recorder, controller := setupTest(t)

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

func TestUserController_NotifyUsers(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	users := []int{1}

	notifyRequest := models.NotifyRequest{
		Users:             users,
		NotificationTitle: "title",
		NotificationText:  "text",
		NotificationType:  "exam_notification",
	}

	jsonBody, _ := json.Marshal(notifyRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	mock.EXPECT().CheckPreference(c.Request.Context(), users[0], notifyRequest.NotificationType).Return(true, nil)
	mock.EXPECT().SendNotifByMobile(c.Request.Context(), users[0], notifyRequest).Return(nil)
	mock.EXPECT().SendNotifByEmail(c.Request.Context(), users[0], notifyRequest).Return(nil)

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestNotifyUsers_InvalidJSON(t *testing.T) {
	_, _, c, recorder, controller := setupTest(t)

	req, _ := http.NewRequest(http.MethodPost, "/users/notify", bytes.NewBufferString("invalid_json"))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid request format")
}

func TestNotifyUsers_CheckPreferenceFalse(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	notifyRequest := models.NotifyRequest{
		Users:             []int{1},
		NotificationTitle: "title",
		NotificationText:  "text",
		NotificationType:  "exam_notification",
	}

	jsonBody, _ := json.Marshal(notifyRequest)
	req := httptest.NewRequest(http.MethodPost, "/users/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mock.EXPECT().CheckPreference(c.Request.Context(), 1, "exam_notification").Return(false, nil)

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestNotifyUsers_PreferenceFalse(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	notifyRequest := models.NotifyRequest{
		Users:             []int{1},
		NotificationTitle: "title",
		NotificationText:  "text",
		NotificationType:  "exam_notification",
	}

	jsonBody, _ := json.Marshal(notifyRequest)
	req := httptest.NewRequest(http.MethodPost, "/users/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mock.EXPECT().CheckPreference(c.Request.Context(), 1, "exam_notification").Return(false, nil)

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestNotifyUsers_SendNotifByMobileError(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	notifyRequest := models.NotifyRequest{
		Users:             []int{1},
		NotificationTitle: "title",
		NotificationText:  "text",
		NotificationType:  "exam_notification",
	}

	jsonBody, _ := json.Marshal(notifyRequest)
	req := httptest.NewRequest(http.MethodPost, "/users/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mock.EXPECT().CheckPreference(c.Request.Context(), 1, "exam_notification").Return(true, nil)
	mock.EXPECT().SendNotifByMobile(c.Request.Context(), 1, notifyRequest).Return(errors.New("mobile error"))

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "mobile error")
}

func TestNotifyUsers_SendNotifByEmailError(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	notifyRequest := models.NotifyRequest{
		Users:             []int{1},
		NotificationTitle: "title",
		NotificationText:  "text",
		NotificationType:  "exam_notification",
	}

	jsonBody, _ := json.Marshal(notifyRequest)
	req := httptest.NewRequest(http.MethodPost, "/users/notify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mock.EXPECT().CheckPreference(c.Request.Context(), 1, "exam_notification").Return(true, nil)
	mock.EXPECT().SendNotifByMobile(c.Request.Context(), 1, notifyRequest).Return(nil)
	mock.EXPECT().SendNotifByEmail(c.Request.Context(), 1, notifyRequest).Return(errors.New("email error"))

	controller.NotifyUsers(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "email error")
}

func TestUserController_SetUserNotifications(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	token := "notification token"

	setupRequest := models.NotificationSetUpRequest{Token: token}

	jsonBody, _ := json.Marshal(setupRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users/1/notifications", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mock.EXPECT().AddNotificationToken(c.Request.Context(), 1, token).Return(nil)

	controller.SetUserNotifications(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_GetUserNotifications(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	req, _ := http.NewRequest(http.MethodGet, "/users/1/notifications", nil)

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req
	token := models.NotificationToken{NotificationToken: "notification token"}
	tokens := models.NotificationTokens{NotificationTokens: []models.NotificationToken{token}}

	mock.EXPECT().GetUserNotificationsToken(c.Request.Context(), 1).Return(tokens, nil)

	controller.GetUserNotifications(c)

	var response models.NotificationTokens

	assert.Equal(t, http.StatusOK, recorder.Code)
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tokens, response)
}

func TestUserController_PasswordReset(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	email := "test@email.com"

	passwordResetRequest := models.PasswordResetRequest{Email: email}
	user := models.User{
		Id:    1,
		Email: email,
	}

	jsonBody, _ := json.Marshal(passwordResetRequest)
	req, _ := http.NewRequest(http.MethodPost, "/users/reset/password", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	c.Request = req

	mock.EXPECT().GetUserByEmail(c.Request.Context(), email).Return(&user, nil)
	mock.EXPECT().StartPasswordReset(c.Request.Context(), user.Id, user.Email).Return(nil)

	controller.PasswordReset(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_AddRule(t *testing.T) {
	_, mockRulesService, c, recorder, controller := setupTest(t)

	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	request := models.Rule{
		Id:                   userId,
		Title:                "title test",
		Description:          "description test",
		EffectiveDate:        time.Time{},
		ApplicationCondition: "condition",
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/rules", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})
	c.Request = req

	mockRulesService.EXPECT().CreateRule(c, request, userId).Return(nil)

	controller.AddRule(c)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_DeleteRule(t *testing.T) {
	_, mockRulesService, c, recorder, controller := setupTest(t)

	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	ruleId := 1

	req := httptest.NewRequest(http.MethodDelete, "/rules/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockRulesService.EXPECT().DeleteRule(c, ruleId, userId).Return(nil)

	controller.DeleteRule(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_GetRules(t *testing.T) {
	_, mockRulesService, c, recorder, controller := setupTest(t)

	req, _ := http.NewRequest(http.MethodGet, "/rules", nil)

	c.Request = req

	rule := models.Rule{
		Id:                   1,
		Title:                "title",
		Description:          "description",
		EffectiveDate:        time.Time{},
		ApplicationCondition: "condition",
	}

	mockRulesService.EXPECT().GetRules(c).Return([]models.Rule{rule}, nil)

	controller.GetRules(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"id":                   float64(rule.Id),
				"Title":                rule.Title,
				"Description":          rule.Description,
				"effectiveDate":        rule.EffectiveDate.Format(time.RFC3339Nano), // if formatted as string
				"ApplicationCondition": rule.ApplicationCondition,
			},
		},
	}
	assert.Equal(t, expected, response)
}

func TestUserController_GetAudits(t *testing.T) {
	_, mockRulesService, c, recorder, controller := setupTest(t)

	req, _ := http.NewRequest(http.MethodGet, "/rules/audit", nil)

	c.Request = req

	audit := models.Audit{
		Id:                   2,
		RuleId:               sql.NullInt64{1, true},
		UserId:               sql.NullInt64{1, true},
		ModificationDate:     time.Time{},
		NatureOfModification: "modification",
	}

	mockRulesService.EXPECT().GetAudits(c).Return([]models.Audit{audit}, nil)

	controller.GetAudits(c)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string][]models.Audit
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	expected := []models.Audit{audit}
	assert.Equal(t, expected, response["data"])
}

func TestUserController_ModifyRule(t *testing.T) {
	_, mockRulesService, c, recorder, controller := setupTest(t)

	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	ruleId := 1

	request := models.RuleModify{
		Title:                "title",
		Description:          "description",
		ApplicationCondition: "condition",
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPut, "/rules/1", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mockRulesService.EXPECT().ModifyRule(c.Request.Context(), ruleId, request, userId).Return(nil)

	controller.ModifyRule(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_ModifyNotifPreference(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	userId := 1

	request := models.NotificationPreferenceRequest{
		NotificationType:       "exam_notification",
		NotificationPreference: false,
	}

	jsonBody, _ := json.Marshal(request)
	req, _ := http.NewRequest(http.MethodPut, "/users/1/notifications/preference", bytes.NewBuffer(jsonBody))

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	mock.EXPECT().SetNotificationPreference(c.Request.Context(), userId, request).Return(nil)

	controller.ModifyNotifPreference(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "null", recorder.Body.String())
}

func TestUserController_GetNotifPreferences(t *testing.T) {
	mock, _, c, recorder, controller := setupTest(t)

	userId := 1

	req, _ := http.NewRequest(http.MethodGet, "/users/1/notifications/preference", nil)

	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	responseExpected := models.NotificationPreference{
		ExamNotification:     true,
		HomeworkNotification: false,
		SocialNotification:   true,
	}

	mock.EXPECT().GetNotificationPreference(c.Request.Context(), userId).Return(&responseExpected, nil)

	controller.GetNotifPreferences(c)

	response := models.NotificationPreference{}

	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, responseExpected, response)
}

func TestUserController_PasswordResetRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, _, c, recorder, controller := setupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/password-reset-redirect?token=abc123", nil)
	c.Request = req

	controller.PasswordResetRedirect(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "myapp://reset-password?token=abc123")
	assert.Contains(t, recorder.Body.String(), "<html>")
}
