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
	GetLatestMessage(context context.Context, userId int) (*models.ChatMessage, error)
}

type chatRepository struct {
	DB *sql.DB
}

func CreateChatsRepo(db *sql.DB) *chatRepository {
	return &chatRepository{DB: db}
}

func (db chatRepository) NewMessage(ctx context.Context, id int, message string, sender string) error {
	_, err := db.DB.ExecContext(ctx, "INSERT INTO ai_chat (user_id, sender, message, rating, feedback) VALUES ($1,$2,$3, 0, $4)", id, sender, message, "")
	if err != nil {
		return err
	}
	return nil
}

func (db chatRepository) GetMessages(ctx context.Context, id int, date string) ([]models.ChatMessage, error) {
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, user_id, sender, message, time_sent, rating, feedback
		FROM ai_chat WHERE user_id = $1 AND DATE(time_sent) >= $2`, id, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		err := rows.Scan(
			&message.MessageId, &message.UserId, &message.Sender, &message.Message, &message.TimeSent, &message.Rating, &message.Feedback,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (db chatRepository) UpdateMessageRating(context context.Context, userId int, messageId int, rating int) error {
	_, err := db.DB.ExecContext(context, "UPDATE ai_chat SET rating = $2 where id = $1 AND user_id = $3", messageId, rating, userId)
	return err
}

func (db chatRepository) UpdateMessageFeedback(context context.Context, userId int, messageId int, feedback string) error {
	_, err := db.DB.ExecContext(context, "UPDATE ai_chat SET feedback = $2 where id = $1 AND user_id = $3", messageId, feedback, userId)
	return err
}

func (db chatRepository) GetLatestMessage(context context.Context, userId int) (*models.ChatMessage, error) {
	query := `
        SELECT id, user_id, sender, message, time_sent, rating, feedback
        FROM ai_chat
        WHERE user_id = $1
        ORDER BY time_sent DESC
        LIMIT 1;
    `
	row := db.DB.QueryRowContext(context, query, userId)

	var chat models.ChatMessage
	err := row.Scan(&chat.MessageId, &chat.UserId, &chat.Sender, &chat.Message, &chat.TimeSent, &chat.Rating, &chat.Feedback)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}
