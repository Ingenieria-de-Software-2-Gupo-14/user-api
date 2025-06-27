package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAuthController(t *testing.T) {
	mockUserService := services.NewMockUserService(t)
	mockLoginAttemptService := services.NewMockLoginAttemptService(t)
	mockVerificationService := services.NewMockVerificationService(t)

	controller := NewAuthController(mockUserService, mockLoginAttemptService, mockVerificationService)
	assert.NotNil(t, controller)
}

func setupTestAuth(t *testing.T) (*services.MockUserService, *services.MockLoginAttemptService, *services.MockVerificationService, *gin.Context, *httptest.ResponseRecorder, *AuthController) {
	gin.SetMode(gin.TestMode)
	mockUserService := services.NewMockUserService(t)
	mockLoginAttemptService := services.NewMockLoginAttemptService(t)
	mockVerificationService := services.NewMockVerificationService(t)

	controller := NewAuthController(mockUserService, mockLoginAttemptService, mockVerificationService)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	return mockUserService, mockLoginAttemptService, mockVerificationService, c, recorder, controller
}

func setupIntegrationTestAuth(db *sql.DB, t *testing.T) (*gin.Context, *httptest.ResponseRecorder, *AuthController, *services.MockEmailSender) {
	gin.SetMode(gin.TestMode)
	email := services.NewMockEmailSender(t)
	repoBlocked := repo.NewBlockedUserRepository(db)
	userService := services.NewUserService(repo.CreateUserRepo(db), repoBlocked, email)
	loginAttemptService := services.NewLoginAttemptService(repo.NewLoginAttemptRepository(db), repoBlocked)
	verificationService := services.NewVerificationService(repo.CreateVerificationRepo(db), email)

	controller := NewAuthController(userService, loginAttemptService, verificationService)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	return c, recorder, controller, email
}

func TestLogin_IntegrationTest(t *testing.T) {
	db, mock, _ := sqlmock.New()
	c, recorder, controller, _ := setupIntegrationTestAuth(db, t)

	password := "testsPassword"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := &models.User{
		Id:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		Name:     "Test User",
		Role:     "user",
		Blocked:  false,
		Verified: true,
	}

	request := models.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("User-Agent", "Mozilla/5.0")

	c.Request = req

	mock.ExpectQuery(`SELECT\s+u\.id,\s*u\.name`).
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "surname", "password", "email", "location", "role", "verified",
			"profile_photo", "description", "created_at", "updated_at", "blocked"}).
			AddRow(user.Id, user.Name, user.Name, user.Password, user.Email, user.Location, user.Role, user.Verified,
				user.ProfilePhoto, user.Description, user.CreatedAt, user.UpdatedAt, user.Blocked))

	mock.ExpectExec(`INSERT INTO login_attempts \(user_id, ip_address, user_agent, successful, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(user.Id, "127.0.0.1", "Mozilla/5.0", true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	controller.Login(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}

func TestResendPin_IntegrationTest(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	c, recorder, controller, mockEmail := setupIntegrationTestAuth(db, t)

	password := "testsPassword"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := &models.User{
		Id:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		Name:     "Test User",
		Role:     "user",
		Blocked:  false,
		Verified: true,
	}

	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend?email="+"test@example.com", nil)
	c.Request = req

	mockDb.ExpectQuery(`SELECT\s+u\.id,\s*u\.name`).
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "surname", "password", "email", "location", "role", "verified",
			"profile_photo", "description", "created_at", "updated_at", "blocked"}).
			AddRow(user.Id, user.Name, user.Name, user.Password, user.Email, user.Location, user.Role, user.Verified,
				user.ProfilePhoto, user.Description, user.CreatedAt, user.UpdatedAt, user.Blocked))

	now := time.Now()
	expected := &models.UserVerification{
		Id:              1,
		UserId:          user.Id,
		UserEmail:       user.Email,
		VerificationPin: "123456",
		PinExpiration:   now.Add(10 * time.Minute),
		CreatedAt:       now,
	}

	mockDb.ExpectQuery(`SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE user_email ILIKE \$1`).
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "user_email", "verification_pin", "pin_expiration", "created_at",
		}).AddRow(
			expected.Id,
			expected.UserId,
			expected.UserEmail,
			expected.VerificationPin,
			expected.PinExpiration,
			expected.CreatedAt,
		))

	mockDb.ExpectExec(`UPDATE verification SET verification_pin = \$2, pin_expiration = \$3 WHERE id = \$1`).
		WithArgs(expected.Id, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	controller.ResendPin(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `null`, recorder.Body.String())
}

func TestRegister(t *testing.T) {
	mockUserService, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	request := models.CreateUserRequest{
		Email:    "test@test.com",
		Password: "test1234",
		Name:     "test",
		Surname:  "test",
		Role:     "student",
		Verified: false,
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctx := req.Context()

	mockUserService.
		EXPECT().
		GetUserByEmail(ctx, request.Email).
		Return(nil, repo.ErrNotFound).
		Once()

	mockUserService.
		EXPECT().
		CreateUser(ctx, request).
		Return(2, nil).
		Once()

	mockVerificationService.
		EXPECT().
		SendVerificationEmail(ctx, 2, request.Email).
		Return(nil).
		Once()

	controller.Register(c)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]int
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response["id"])
}

func TestRegister_StatusBadRequest_Error(t *testing.T) {
	userService, _, verificationService, c, w, controller := setupTestAuth(t)

	c.Request = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	userService.AssertExpectations(t)
	verificationService.AssertExpectations(t)
}

func TestRegister_ExistingUserVerified_Error(t *testing.T) {
	userService, _, verificationService, c, w, controller := setupTestAuth(t)

	request := models.CreateUserRequest{
		Email:    "test@test.com",
		Password: "test1234",
		Name:     "test",
		Surname:  "test",
		Role:     "student",
		Verified: true,
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	userService.EXPECT().GetUserByEmail(c.Request.Context(), "test@test.com").Return(&models.User{Verified: true}, nil)

	controller.Register(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	userService.AssertExpectations(t)
	verificationService.AssertExpectations(t)
}

func TestRegisterAdmin(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	request := models.CreateUserRequest{
		Email:    "test@test.com",
		Password: "test1234",
		Name:     "test",
		Surname:  "test",
		Role:     "admin",
		Verified: true,
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctx := req.Context()

	mockUserService.
		EXPECT().
		GetUserByEmail(ctx, request.Email).
		Return(nil, repo.ErrNotFound).
		Once()

	mockUserService.
		EXPECT().
		CreateUser(ctx, request).
		Return(2, nil).
		Once()

	controller.RegisterAdmin(c)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response map[string]int
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response["id"])
}

func TestVerifyRegistration(t *testing.T) {
	mockUserService, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	pinId := 01
	pinCode := "123456"
	fullPin := fmt.Sprintf("%d-%s", pinId, pinCode)
	userID := 1
	email := "test@example.com"

	expectedVerification := &models.UserVerification{
		Id:              pinId,
		UserId:          userID,
		UserEmail:       email,
		VerificationPin: pinCode,
		PinExpiration:   time.Now().Add(5 * time.Minute), // not expired
	}

	request := models.EmailVerifiaction{VerificationPin: fullPin}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	ctx := req.Context()

	mockVerificationService.
		EXPECT().
		GetVerification(ctx, pinId).
		Return(expectedVerification, nil).
		Once()

	mockUserService.
		EXPECT().
		VerifyUser(ctx, userID).
		Return(nil).
		Once()

	mockVerificationService.
		EXPECT().
		DeleteByUserId(ctx, userID).
		Return(nil).
		Once()

	controller.VerifyRegistration(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User verified and created successfully", response["message"])
}

func TestVerifyRegistration_InvalidJSON(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.VerifyRegistration(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestVerifyRegistration_InvalidPinFormat(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	request := models.EmailVerifiaction{VerificationPin: "invalidpin"}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.VerifyRegistration(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestVerifyRegistration_InvalidPinId(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	request := models.EmailVerifiaction{VerificationPin: "abc-123456"}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.VerifyRegistration(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestVerifyRegistration_PinExpired(t *testing.T) {
	_, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	expired := &models.UserVerification{
		Id: 1, UserId: 2, VerificationPin: "123456", PinExpiration: time.Now().Add(-time.Minute),
	}
	request := models.EmailVerifiaction{VerificationPin: "1-123456"}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()

	c.Request = req.WithContext(ctx)

	mockVerificationService.EXPECT().GetVerification(ctx, 1).Return(expired, nil).Once()

	controller.VerifyRegistration(c)
	assert.Equal(t, http.StatusGone, recorder.Code)
}

func TestVerifyRegistration_PinMismatch(t *testing.T) {
	_, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	wrongPin := &models.UserVerification{
		Id: 1, UserId: 2, VerificationPin: "1", PinExpiration: time.Now().Add(time.Minute),
	}
	request := models.EmailVerifiaction{VerificationPin: "1-123456"}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/users/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := req.Context()

	c.Request = req.WithContext(ctx)

	mockVerificationService.EXPECT().GetVerification(ctx, 1).Return(wrongPin, nil).Once()

	controller.VerifyRegistration(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestLogin(t *testing.T) {
	mockUserService, mockLoginService, _, c, recorder, controller := setupTestAuth(t)

	password := "testsPassword"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := &models.User{
		Id:       1,
		Email:    "test@example.com",
		Password: hashedPassword,
		Name:     "Test User",
		Role:     "user",
		Blocked:  false,
		Verified: true,
	}

	request := models.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1"
	req.Header.Set("User-Agent", "Mozilla/5.0")

	c.Request = req

	ctx := req.Context()

	mockUserService.
		EXPECT().
		GetUserByEmail(ctx, user.Email).
		Return(user, nil).
		Once()

	mockLoginService.
		EXPECT().
		AddLoginAttempt(c, user.Id, "127.0.0.1", "Mozilla/5.0", true).
		Return(nil).
		Once()

	controller.Login(c)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["token"])
}

func TestLogin_BadRequest(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	request := models.LoginRequest{
		Email:    "a",
		Password: "password",
	}

	jsonBody, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.Login(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
func TestLogin_UserNotFound(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	request := models.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockUserService.EXPECT().
		GetUserByEmail(c.Request.Context(), request.Email).
		Return(nil, sql.ErrNoRows).
		Once()

	controller.Login(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockUserService, mockLoginService, _, c, recorder, controller := setupTestAuth(t)

	hashedPassword, _ := utils.HashPassword("correctPassword")
	user := &models.User{
		Id:       1,
		Email:    "test@test.com",
		Password: hashedPassword,
		Verified: true,
	}

	request := models.LoginRequest{
		Email:    user.Email,
		Password: "wrongPassword",
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.RemoteAddr = "127.0.0.1"
	c.Request = req

	mockUserService.EXPECT().
		GetUserByEmail(c.Request.Context(), user.Email).
		Return(user, nil).
		Once()

	mockLoginService.EXPECT().
		AddLoginAttempt(c, user.Id, "127.0.0.1", "Mozilla/5.0", false).
		Return(nil).
		Once()

	controller.Login(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestLogin_UserBlocked(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	password := "password"
	hashedPassword, _ := utils.HashPassword(password)

	user := &models.User{
		Id:       1,
		Email:    "test@test.com",
		Password: hashedPassword,
		Blocked:  true,
		Verified: true,
	}

	request := models.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockUserService.EXPECT().
		GetUserByEmail(c.Request.Context(), user.Email).
		Return(user, nil).
		Once()

	controller.Login(c)

	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestLogin_UserNotVerified(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	password := "password"
	hashedPassword, _ := utils.HashPassword(password)

	user := &models.User{
		Id:       1,
		Email:    "test@test.com",
		Password: hashedPassword,
		Blocked:  false,
		Verified: false,
	}

	request := models.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	mockUserService.EXPECT().
		GetUserByEmail(c.Request.Context(), user.Email).
		Return(user, nil).
		Once()

	controller.Login(c)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestLogout(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	c.Request = req

	// Act
	controller.Logout(c)

	// Assert
	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, "/", recorder.Header().Get("Location"))

	// Check Set-Cookie header
	setCookie := recorder.Header().Get("Set-Cookie")
	assert.Contains(t, setCookie, "Authorization=")
	assert.Contains(t, setCookie, "Max-Age=0") // or a negative expiration
}

func TestResendPin(t *testing.T) {
	mockUserService, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	email := "test@example.com"
	user := &models.User{
		Id:    1,
		Email: email,
		Name:  "Test",
	}

	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend?email="+email, nil)
	c.Request = req

	mockUserService.
		EXPECT().
		GetUserByEmail(c.Request.Context(), email).
		Return(user, nil).
		Once()

	mockVerificationService.
		EXPECT().
		UpdatePin(c.Request.Context(), user.Id, email).
		Return(nil).
		Once()

	controller.ResendPin(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `null`, recorder.Body.String())
}

func TestResendPin_MissingEmail(t *testing.T) {
	_, _, _, c, recorder, controller := setupTestAuth(t)

	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend", nil) // no query param
	c.Request = req

	controller.ResendPin(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Email is required")
}

func TestResendPin_UserNotFound(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	email := "nonexistent@example.com"
	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend?email="+email, nil)
	c.Request = req

	mockUserService.
		EXPECT().
		GetUserByEmail(c.Request.Context(), email).
		Return(nil, repo.ErrNotFound).
		Once()

	controller.ResendPin(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid verification")
}

func TestResendPin_GetUserByEmailInternalError(t *testing.T) {
	mockUserService, _, _, c, recorder, controller := setupTestAuth(t)

	email := "test@example.com"
	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend?email="+email, nil)
	c.Request = req

	mockUserService.
		EXPECT().
		GetUserByEmail(c.Request.Context(), email).
		Return(nil, errors.New("db error")).
		Once()

	controller.ResendPin(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "db error")
}

func TestResendPin_UpdatePinError(t *testing.T) {
	mockUserService, _, mockVerificationService, c, recorder, controller := setupTestAuth(t)

	email := "test@example.com"
	user := &models.User{Id: 1, Email: email}

	req := httptest.NewRequest(http.MethodPost, "/users/verify/resend?email="+email, nil)
	c.Request = req

	mockUserService.
		EXPECT().
		GetUserByEmail(c.Request.Context(), email).
		Return(user, nil).
		Once()

	mockVerificationService.
		EXPECT().
		UpdatePin(c.Request.Context(), user.Id, email).
		Return(errors.New("send error")).
		Once()

	controller.ResendPin(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "send error")
}
