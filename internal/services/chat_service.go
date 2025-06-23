package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ChatService interface {
	UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error
	UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error
	NewUserMessage(context context.Context, id int, message models.ChatMessageRequest) error
	GetMessages(ctx context.Context, id int) ([]models.ChatMessage, error)
	SendToAi(ctx context.Context, id int, message models.ChatMessageRequest) (*models.ChatMessage, error)
}

type chatService struct {
	chatRepo repo.ChatRepository
}

func NewChatsService(chatRepo repo.ChatRepository) *chatService {
	return &chatService{chatRepo: chatRepo}
}

func (s *chatService) UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error {
	return s.chatRepo.UpdateMessageRating(context, userId, messageId, rating)
}

func (s *chatService) UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error {
	return s.chatRepo.UpdateMessageFeedback(context, userId, messageId, feedback)
}

func (s *chatService) NewUserMessage(context context.Context, id int, message models.ChatMessageRequest) error {
	return s.chatRepo.NewMessage(context, id, message.Message, "user")
}

func (s *chatService) GetMessages(ctx context.Context, id int) ([]models.ChatMessage, error) {
	return s.chatRepo.GetMessages(ctx, id, time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
}

func (s *chatService) SendToAi(ctx context.Context, id int, message models.ChatMessageRequest) (*models.ChatMessage, error) {
	messages, err := s.GetMessages(ctx, id)
	if err != nil {
		return nil, err
	}
	var messagesToSend []models.Message
	contextPropmpt := models.Message{
		Role:    "system",
		Content: "You are an ai assistant for an app similar to google classroom, named classconnect, where students can enroll into courses, upload their homework, upload their exam results, communicate through a forum with both teachers and students. Teachers can manage courses, their students, check their performance, assign a user as an assistant teacher. as an assistant you should help both teachers and students, and should recommend bibliography or other useful material if needed. while chatting with you the users might rate or give feedback towards your answers, don't forget to take them into account.Finally its better if you help a student arrive to an answer rather than just giving it to them",
	}
	messagesToSend = append(messagesToSend, contextPropmpt)
	for _, dbMessage := range messages {
		var messageToSend models.Message
		if dbMessage.Sender == "ai" {
			dbMessage.Sender = "assistant"
		}
		messageToSend.Role = dbMessage.Sender
		messageToSend.Content = dbMessage.Message
		messagesToSend = append(messagesToSend, messageToSend)
		messagesToSend = checkForFeedbackAndRating(messagesToSend, dbMessage)
	}
	url := "https://api.openai.com/v1/chat/completions"

	requestBody := models.OpenAIRequest{
		Model:    "gpt-3.5-turbo", // or gpt-4o if your key allows
		Messages: messagesToSend,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAT_GPT_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp models.OpenAIResponse
	err = json.Unmarshal(body, &openAIResp)
	if err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) > 0 {
		err = s.chatRepo.NewMessage(ctx, id, openAIResp.Choices[0].Message.Content, "assistant")
		if err != nil {
			return nil, err
		}
		return s.chatRepo.GetLatestMessage(ctx, id)
	}
	return nil, errors.New("Something went wrong with the Ai")
}

func checkForFeedbackAndRating(messagesToSend []models.Message, dbMessage models.ChatMessage) []models.Message {
	if dbMessage.Rating == 0 && dbMessage.Feedback == "" {
		return messagesToSend
	}
	var feedbackRating string
	if dbMessage.Feedback != "" && dbMessage.Sender == "assistant" {
		feedbackRating = feedbackRating + "feedback: " + dbMessage.Feedback + ". "
	}
	if dbMessage.Rating != 0 && dbMessage.Sender == "assistant" {
		feedbackRating = feedbackRating + "rating: " + strconv.Itoa(dbMessage.Rating) + ". "
	}
	var messageToSend models.Message
	messageToSend.Role = "user"
	messageToSend.Content = feedbackRating
	return append(messagesToSend, messageToSend)
}
