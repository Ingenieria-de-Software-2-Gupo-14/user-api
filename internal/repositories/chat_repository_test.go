package repositories

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateChatsRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)
	assert.NotNil(t, repo)
}

func TestChatRepository_NewMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)

	ctx := context.Background()
	userID := 1
	sender := "user"
	message := "test"

	mock.ExpectExec(`INSERT INTO ai_chat \(user_id, sender, message, rating, feedback\) VALUES \(\$1,\$2,\$3, 0, \$4\)`).
		WithArgs(userID, sender, message, "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.NewMessage(ctx, userID, message, sender)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)

	ctx := context.Background()
	userID := 42
	sinceDate := "2024-01-01"

	expectedMessages := []models.ChatMessage{
		{
			MessageId: 1,
			UserId:    userID,
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

	mock.ExpectQuery(`SELECT id, user_id, sender, message, time_sent, rating, feedback FROM ai_chat WHERE user_id = \$1 AND DATE\(time_sent\) >= \$2`).
		WithArgs(userID, sinceDate).
		WillReturnRows(rows)

	messages, err := repo.GetMessages(ctx, userID, sinceDate)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_GetLatestMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)

	ctx := context.Background()
	userID := 1
	sinceDate := "2024-01-01"

	expected := models.ChatMessage{
		MessageId: 1,
		UserId:    userID,
		Sender:    "user",
		Message:   "test",
		TimeSent:  sinceDate,
		Rating:    0,
		Feedback:  "",
	}

	// Expect the query
	mock.ExpectQuery(`SELECT id, user_id, sender, message, time_sent, rating, feedback
        FROM ai_chat
        WHERE user_id = \$1
        ORDER BY time_sent DESC
        LIMIT 1;`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "sender", "message", "time_sent", "rating", "feedback",
		}).
			AddRow(expected.MessageId, expected.UserId, expected.Sender, expected.Message,
				expected.TimeSent, expected.Rating, expected.Feedback),
		)

	msg, err := repo.GetLatestMessage(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, &expected, msg)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_UpdateMessageFeedback(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)

	ctx := context.Background()
	userId := 1
	messageId := 1
	feedback := "feedback"

	// Expect the Exec query with correct args
	mock.ExpectExec(`UPDATE ai_chat SET feedback = \$2 where id = \$1 AND user_id = \$3`).
		WithArgs(messageId, feedback, userId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.UpdateMessageFeedback(ctx, userId, messageId, feedback)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestChatRepository_UpdateMessageRating(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateChatsRepo(db)

	ctx := context.Background()
	userId := 1
	messageId := 1
	rating := 2

	// Expect the Exec query with correct args
	mock.ExpectExec(`UPDATE ai_chat SET rating = \$2 where id = \$1 AND user_id = \$3`).
		WithArgs(messageId, rating, userId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.UpdateMessageRating(ctx, userId, messageId, rating)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
