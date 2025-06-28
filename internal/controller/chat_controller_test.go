package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	s "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/services"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*s.MockChatService, *gin.Context, *httptest.ResponseRecorder, *ChatController) {
	gin.SetMode(gin.TestMode)
	mockService := s.NewMockChatService(t)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	userChat := NewChatsController(mockService)
	return mockService, c, recorder, userChat
}

func TestNewChatsController(t *testing.T) {
	mockService := s.NewMockChatService(t)
	result := NewChatsController(mockService)
	assert.NotNil(t, result)
}

func TestChatController_SendMessage(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	message := "test message"
	chatRequestBody := models.ChatMessageRequest{Message: message}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req

	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	ai_response := models.ChatMessage{
		MessageId: 1,
		UserId:    userId,
		Message:   "ai test message",
		Sender:    "assistant",
		TimeSent:  "01/01/2025",
		Rating:    0,
		Feedback:  "",
	}

	mockService.EXPECT().NewUserMessage(c.Request.Context(), userId, chatRequestBody).Return(nil)
	mockService.EXPECT().SendToAi(c.Request.Context(), userId, chatRequestBody).Return(&ai_response, nil)

	userChat.SendMessage(c)

	assert.Equal(t, http.StatusOK, recorder.Code)
	response := models.ChatMessage{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, ai_response, response)
}

func TestChatController_SendMessage_Bad_Request(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	message := ""
	chatRequestBody := models.ChatMessageRequest{Message: message}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	userChat.SendMessage(c)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusBadRequest,
		Title: http.StatusText(http.StatusBadRequest),
		Error: "Invalid request format",
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_SendMessage_Incorrect_Token(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	message := "test message"
	chatRequestBody := models.ChatMessageRequest{Message: message}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: "a",
	})

	// Assign request to context
	c.Request = req

	userChat.SendMessage(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestChatController_SendMessage_Service_Fails1(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	message := "test message"
	chatRequestBody := models.ChatMessageRequest{Message: message}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	expectedErr := errors.New("database error")

	mockService.EXPECT().NewUserMessage(c.Request.Context(), userId, chatRequestBody).Return(expectedErr)

	userChat.SendMessage(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusInternalServerError,
		Title: http.StatusText(http.StatusInternalServerError),
		Error: expectedErr.Error(),
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_SendMessage_Service_Fails2(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	message := "test message"
	chatRequestBody := models.ChatMessageRequest{Message: message}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	expectedErr := errors.New("database error")

	mockService.EXPECT().NewUserMessage(c.Request.Context(), userId, chatRequestBody).Return(nil)
	mockService.EXPECT().SendToAi(c.Request.Context(), userId, chatRequestBody).Return(nil, expectedErr)

	userChat.SendMessage(c)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusInternalServerError,
		Title: http.StatusText(http.StatusInternalServerError),
		Error: expectedErr.Error(),
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_GetMessages(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	message := models.ChatMessage{
		MessageId: 1,
		UserId:    userId,
		Message:   "test message",
		Sender:    "user",
		TimeSent:  "01/01/2025",
		Rating:    0,
		Feedback:  "",
	}
	expectedMessages := []models.ChatMessage{message}

	mockService.EXPECT().GetMessages(c.Request.Context(), userId).Return(expectedMessages, nil)

	userChat.GetMessages(c)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []models.ChatMessage
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, response)
}

func TestChatController_GetMessages_Incorrect_Token(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: "a",
	})

	// Assign request to context
	c.Request = req

	userChat.GetMessages(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestChatController_GetMessages_Service_Fails(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	expectedErr := errors.New("database error")

	mockService.EXPECT().GetMessages(c.Request.Context(), userId).Return(nil, expectedErr)

	userChat.GetMessages(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusInternalServerError,
		Title: http.StatusText(http.StatusInternalServerError),
		Error: expectedErr.Error(),
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_FeedbackMessage(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	feedback := "feedback"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatFeedbackRequest{Feedback: feedback}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/feedback", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	mockService.EXPECT().UpdateMessageFeedback(c.Request.Context(), userId, messageId, feedback).Return(nil)

	userChat.FeedbackMessage(c)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestChatController_FeedbackMessage_Incorrect_Token(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	messageId := 1
	feedback := "feedback"

	chatRequestBody := models.ChatFeedbackRequest{Feedback: feedback}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/feedback", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: "a",
	})

	// Assign request to context
	c.Request = req
	c.AddParam("message_id", strconv.Itoa(messageId))

	userChat.FeedbackMessage(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestChatController_FeedbackMessage_Incorrect_Path_Variable(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	userId := 1
	//messageId := "failure test"
	email := "test@test.com"
	name := "test"
	role := "user"
	feedback := "feedback"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatFeedbackRequest{Feedback: feedback}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	req := httptest.NewRequest(http.MethodPut, "/chat//feedback", body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	userChat.FeedbackMessage(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusBadRequest,
		Title: http.StatusText(http.StatusBadRequest),
		Error: "Invalid message ID format",
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_FeedbackMessage_Bad_Request(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	feedback := ""
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatFeedbackRequest{Feedback: feedback}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/feedback", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	userChat.FeedbackMessage(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusBadRequest,
		Title: http.StatusText(http.StatusBadRequest),
		Error: "Invalid request format",
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_FeedbackMessage_Service_Fails(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	feedback := "feedback"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatFeedbackRequest{Feedback: feedback}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/feedback", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	expectedErr := errors.New("database error")

	mockService.EXPECT().UpdateMessageFeedback(c.Request.Context(), userId, messageId, feedback).Return(expectedErr)

	userChat.FeedbackMessage(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusInternalServerError,
		Title: http.StatusText(http.StatusInternalServerError),
		Error: expectedErr.Error(),
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_RateMessage(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	rating := 2
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatRatingRequest{Rating: rating}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/rate", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	mockService.EXPECT().UpdateMessageRating(c.Request.Context(), userId, messageId, rating).Return(nil)

	userChat.RateMessage(c)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestChatController_RateMessage_Incorrect_Token(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	messageId := 1
	rating := 2

	chatRequestBody := models.ChatRatingRequest{Rating: rating}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/rate", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: "a",
	})

	// Assign request to context
	c.Request = req
	c.AddParam("message_id", strconv.Itoa(messageId))

	userChat.RateMessage(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestChatController_RateMessage_Incorrect_Path_Variable(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	rating := 2
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatRatingRequest{Rating: rating}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	req := httptest.NewRequest(http.MethodPut, "/chat//rate", body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	userChat.RateMessage(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusBadRequest,
		Title: http.StatusText(http.StatusBadRequest),
		Error: "Invalid message ID format",
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_RateMessage_Bad_Request(t *testing.T) {
	_, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	rating := 0
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatRatingRequest{Rating: rating}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/rate", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	userChat.RateMessage(c)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusBadRequest,
		Title: http.StatusText(http.StatusBadRequest),
		Error: "Invalid request format",
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestChatController_RateMessage_Service_Fails(t *testing.T) {
	mockService, c, recorder, userChat := setupTest(t)
	userId := 1
	messageId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	rating := 2
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	chatRequestBody := models.ChatRatingRequest{Rating: rating}
	jsonValue, err := json.Marshal(chatRequestBody)
	assert.NoError(t, err)
	body := bytes.NewBuffer(jsonValue)

	urlTarget := fmt.Sprintf("/chat/%d/rate", messageId)
	req := httptest.NewRequest(http.MethodPut, urlTarget, body)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})
	c.AddParam("message_id", strconv.Itoa(messageId))

	expectedErr := errors.New("database error")

	mockService.EXPECT().UpdateMessageRating(c.Request.Context(), userId, messageId, rating).Return(expectedErr)

	userChat.RateMessage(c)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	expectedResponse := utils.HTTPError{
		Code:  http.StatusInternalServerError,
		Title: http.StatusText(http.StatusInternalServerError),
		Error: expectedErr.Error(),
	}

	response := utils.HTTPError{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

// integration tests

func setupIntegrationTest(db *sql.DB, t *testing.T) (*gin.Context, *httptest.ResponseRecorder, *ChatController) {
	gin.SetMode(gin.TestMode)
	service := s.NewChatsService(repositories.CreateChatsRepo(db))
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	userChat := NewChatsController(service)
	return c, recorder, userChat
}

func TestChatController_GetMessages_IntegrationTest(t *testing.T) {
	db, mockDb, _ := sqlmock.New()
	c, recorder, userChat := setupIntegrationTest(db, t)
	userId := 1
	email := "test@test.com"
	name := "test"
	role := "user"
	token, err := models.GenerateToken(userId, email, name, role)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	req.AddCookie(&http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	// Assign request to context
	c.Request = req
	c.Set("claims", &models.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: strconv.Itoa(userId),
		},
		Email: email,
		Name:  name,
		Role:  role,
	})

	sinceDate := "2024-01-01"

	expectedMessages := []models.ChatMessage{
		{
			MessageId: 1,
			UserId:    userId,
			Sender:    "user",
			Message:   "test",
			TimeSent:  sinceDate,
			Rating:    0,
			Feedback:  "",
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "sender", "message", "time_sent", "rating", "feedback",
	}).
		AddRow(
			expectedMessages[0].MessageId,
			expectedMessages[0].UserId,
			expectedMessages[0].Sender,
			expectedMessages[0].Message,
			expectedMessages[0].TimeSent,
			expectedMessages[0].Rating,
			expectedMessages[0].Feedback,
		)

	mockDb.ExpectQuery(`SELECT id, user_id, sender, message, time_sent, rating, feedback FROM ai_chat WHERE user_id = \$1 AND DATE\(time_sent\) >= \$2`).
		WithArgs(userId, time.Now().AddDate(0, 0, -1).Format("2006-01-02")).
		WillReturnRows(rows)

	userChat.GetMessages(c)
	assert.Equal(t, http.StatusOK, recorder.Code)

	var response []models.ChatMessage
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, response)
}
