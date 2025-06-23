package services

import (
	"context"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

func TestNewChatsService(t *testing.T) {
	mockRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(mockRepo)

	assert.NotNil(t, service)
}

func TestUpdateMessageRating(t *testing.T) {
	mockRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(mockRepo)

	ctx := context.Background()
	userId := 1
	messageId := 2
	rating := 5

	mockRepo.
		EXPECT().
		UpdateMessageRating(ctx, userId, messageId, rating).
		Return(nil)

	err := service.UpdateMessageRating(ctx, userId, messageId, rating)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMessageFeedback(t *testing.T) {
	mockRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(mockRepo)

	ctx := context.Background()
	userId := 1
	messageId := 2
	feedback := "Test Feedback"

	mockRepo.
		EXPECT().
		UpdateMessageFeedback(ctx, userId, messageId, feedback).
		Return(nil)

	err := service.UpdateMessageFeedback(ctx, userId, messageId, feedback)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestNewUserMessage(t *testing.T) {
	mockRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(mockRepo)

	ctx := context.Background()
	chatId := 1
	req := models.ChatMessageRequest{Message: "Test message"}

	mockRepo.
		EXPECT().
		NewMessage(ctx, chatId, req.Message, "user").
		Return(nil)

	err := service.NewUserMessage(ctx, chatId, req)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetMessages(t *testing.T) {
	mockRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(mockRepo)

	ctx := context.Background()
	chatId := 1
	expected := []models.ChatMessage{
		{Sender: "user", Message: "Test Message"},
		{Sender: "ai", Message: "Test Response"},
	}

	// Calculate expected date
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	mockRepo.
		EXPECT().
		GetMessages(ctx, chatId, date).
		Return(expected, nil)

	messages, err := service.GetMessages(ctx, chatId)
	assert.NoError(t, err)
	assert.Equal(t, expected, messages)

	mockRepo.AssertExpectations(t)
}

func TestSendToAi(t *testing.T) {
	ctx := context.Background()
	chatRepo := repo.NewMockChatRepository(t)
	service := NewChatsService(chatRepo)

	userId := 1
	req := models.ChatMessageRequest{
		Message: "test",
	}

	// Mock chat messages
	chatRepo.EXPECT().
		GetMessages(ctx, userId, mock.Anything).
		Return([]models.ChatMessage{
			{Sender: "user", Message: "Hello"},
			{Sender: "ai", Message: "Hi, how can I help?", Rating: 3, Feedback: "test feedback"},
		}, nil)

	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.openai.com/v1/chat/completions",
		func(req *http.Request) (*http.Response, error) {
			response := `{
                "choices": [{
                    "message": {
                        "role": "assistant",
                        "content": "Sure! Here’s a summary..."
                    }
                }]
            }`
			return httpmock.NewStringResponse(200, response), nil
		})

	// Expect new message saved
	chatRepo.EXPECT().
		NewMessage(ctx, userId, "Sure! Here’s a summary...", "assistant").
		Return(nil)

	// Expect latest message returned
	expected := &models.ChatMessage{
		MessageId: 1,
		UserId:    1,
		Message:   "test",
		Sender:    "user",
		TimeSent:  "20/03/2025",
		Rating:    2,
		Feedback:  "test feedback",
	}
	chatRepo.EXPECT().
		GetLatestMessage(ctx, userId).
		Return(expected, nil)

	result, err := service.SendToAi(ctx, userId, req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
