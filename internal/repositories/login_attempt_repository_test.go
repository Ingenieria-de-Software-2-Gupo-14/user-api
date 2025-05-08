package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewLoginAttemptRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewLoginAttemptRepository(db)
	assert.NotNil(t, repo)
}

func TestLoginAttemptDB_AddLoginAttempt(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	successful := true

	mock.ExpectExec(`INSERT INTO login_attempts \(user_id, ip_address, user_agent, successful, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(userID, ipAddress, userAgent, successful, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	err = repo.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptDB_BadLoginAttemptsInLast10Minutes_NoAttempts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	expectedCount := 0

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM login_attempts WHERE user_id = \$1 AND successful = false AND created_at >= NOW\(\) - INTERVAL '10 minutes'`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	count, err := repo.BadLoginAttemptsInLast10Minutes(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptDB_BadLoginAttemptsInLast10Minutes_WithAttempts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	expectedCount := 3

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM login_attempts WHERE user_id = \$1 AND successful = false AND created_at >= NOW\(\) - INTERVAL '10 minutes'`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	count, err := repo.BadLoginAttemptsInLast10Minutes(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptDB_GetLoginsByUserId_NoAttempts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT id, user_id, ip_address, user_agent, successful, created_at FROM login_attempts WHERE user_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(userID, limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "ip_address", "user_agent", "successful", "created_at"}))

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	attempts, err := repo.GetLoginsByUserId(ctx, userID, limit, offset)
	assert.NoError(t, err)
	assert.Empty(t, attempts)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptDB_GetLoginsByUserId_WithAttempts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	limit := 10
	offset := 0
	now := time.Now()

	// Create test data
	attempt1 := &models.LoginAttempt{
		ID:         1,
		UserID:     userID,
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Successful: true,
		CreatedAt:  now.Add(-1 * time.Hour),
	}

	attempt2 := &models.LoginAttempt{
		ID:         2,
		UserID:     userID,
		IPAddress:  "192.168.1.2",
		UserAgent:  "Chrome/90.0",
		Successful: false,
		CreatedAt:  now,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "ip_address", "user_agent", "successful", "created_at"}).
		AddRow(attempt2.ID, attempt2.UserID, attempt2.IPAddress, attempt2.UserAgent, attempt2.Successful, attempt2.CreatedAt).
		AddRow(attempt1.ID, attempt1.UserID, attempt1.IPAddress, attempt1.UserAgent, attempt1.Successful, attempt1.CreatedAt)

	mock.ExpectQuery(`SELECT id, user_id, ip_address, user_agent, successful, created_at FROM login_attempts WHERE user_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	attempts, err := repo.GetLoginsByUserId(ctx, userID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, attempts, 2)

	// Check first attempt (most recent due to ORDER BY created_at DESC)
	assert.Equal(t, attempt2.ID, attempts[0].ID)
	assert.Equal(t, attempt2.UserID, attempts[0].UserID)
	assert.Equal(t, attempt2.IPAddress, attempts[0].IPAddress)
	assert.Equal(t, attempt2.UserAgent, attempts[0].UserAgent)
	assert.Equal(t, attempt2.Successful, attempts[0].Successful)
	assert.Equal(t, attempt2.CreatedAt.Unix(), attempts[0].CreatedAt.Unix())

	// Check second attempt
	assert.Equal(t, attempt1.ID, attempts[1].ID)
	assert.Equal(t, attempt1.UserID, attempts[1].UserID)
	assert.Equal(t, attempt1.IPAddress, attempts[1].IPAddress)
	assert.Equal(t, attempt1.UserAgent, attempts[1].UserAgent)
	assert.Equal(t, attempt1.Successful, attempts[1].Successful)
	assert.Equal(t, attempt1.CreatedAt.Unix(), attempts[1].CreatedAt.Unix())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptDB_GetLoginsByUserId_WithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	limit := 1
	offset := 1
	now := time.Now()

	// Second attempt will be returned due to offset
	attempt2 := &models.LoginAttempt{
		ID:         2,
		UserID:     userID,
		IPAddress:  "192.168.1.2",
		UserAgent:  "Chrome/90.0",
		Successful: false,
		CreatedAt:  now.Add(-1 * time.Hour),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "ip_address", "user_agent", "successful", "created_at"}).
		AddRow(attempt2.ID, attempt2.UserID, attempt2.IPAddress, attempt2.UserAgent, attempt2.Successful, attempt2.CreatedAt)

	mock.ExpectQuery(`SELECT id, user_id, ip_address, user_agent, successful, created_at FROM login_attempts WHERE user_id = \$1 ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)

	ctx := context.Background()
	repo := NewLoginAttemptRepository(db)

	attempts, err := repo.GetLoginsByUserId(ctx, userID, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, attempts, 1)

	// Check the attempt
	assert.Equal(t, attempt2.ID, attempts[0].ID)
	assert.Equal(t, attempt2.UserID, attempts[0].UserID)
	assert.Equal(t, attempt2.IPAddress, attempts[0].IPAddress)
	assert.Equal(t, attempt2.UserAgent, attempts[0].UserAgent)
	assert.Equal(t, attempt2.Successful, attempts[0].Successful)
	assert.Equal(t, attempt2.CreatedAt.Unix(), attempts[0].CreatedAt.Unix())

	assert.NoError(t, mock.ExpectationsWereMet())
}
