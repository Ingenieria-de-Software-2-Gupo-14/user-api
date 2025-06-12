package services

import (
	"context"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"time"
)

type ChatService interface {
	UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error
	UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error
	NewUserMessage(context context.Context, id int, message models.ChatMessageRequest) error
	GetMessages(ctx context.Context, id int) ([]models.ChatMessage, error)
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
