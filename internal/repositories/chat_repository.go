package repositories

import (
	"context"
	"database/sql"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	_ "github.com/lib/pq"
)

type ChatRepository interface {
	UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error
	UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error
	NewMessage(ctx context.Context, id int, message string, sender string) error
	GetMessages(ctx context.Context, id int, date string) ([]models.ChatMessage, error)
}

type chatRepository struct {
	DB *sql.DB
}

func CreateChatsRepo(db *sql.DB) *chatRepository {
	return &chatRepository{DB: db}
}

func (db chatRepository) NewMessage(ctx context.Context, id int, message string, sender string) error {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO ai_chat (user_id, sender, message) VALUES ($1,$2,$3)", id, sender, message)
	if err != nil {
		return err
	}
	return nil
}

func (db chatRepository) GetMessages(ctx context.Context, id int, date string) ([]models.ChatMessage, error) {
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, user_id, sender, message, time_sent, rating, feedback
		FROM ai_chat WHERE user_id = $1, DATE(time_sent) >= $2`, id, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		// Scan sin bad_login_attempts ni blocked
		err := rows.Scan(
			&message.MessageId, &message.UserId, &message.Sender, &message.TimeSent, &message.Rating, &message.Feedback,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (db chatRepository) UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error {
	_, err := db.DB.ExecContext(context, "UPDATE ai_chat SET rating = $2 where id = $1", messageId, rating)
	return err
}

func (db chatRepository) UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error {
	_, err := db.DB.ExecContext(context, "UPDATE ai_chat SET feedback = $2 where id = $1", messageId, feedback)
	return err
}
