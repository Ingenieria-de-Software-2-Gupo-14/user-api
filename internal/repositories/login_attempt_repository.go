package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

type LoginAttemptRepository interface {
	AddLoginAttempt(ctx context.Context, userID int, ipAddress, userAgent string, successful bool) error
	BadLoginAttemptsInLast10Minutes(ctx context.Context, userID int) (int, error)
	GetLoginsByUserId(ctx context.Context, userID int, limit, offset int) ([]*models.LoginAttempt, error)
}

type LoginAttemptDB struct {
	DB *sql.DB
}

func NewLoginAttemptRepository(db *sql.DB) *LoginAttemptDB {
	return &LoginAttemptDB{DB: db}
}

func (db *LoginAttemptDB) AddLoginAttempt(ctx context.Context, userID int, ipAddress, userAgent string, successful bool) error {
	query := `
		INSERT INTO login_attempts (user_id, ip_address, user_agent, successful, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.DB.ExecContext(ctx, query, userID, ipAddress, userAgent, successful, time.Now())
	return err
}

func (db *LoginAttemptDB) BadLoginAttemptsInLast10Minutes(ctx context.Context, userID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM login_attempts
		WHERE user_id = $1 AND successful = false AND created_at >= NOW() - INTERVAL '10 minutes'`
	var count int
	err := db.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *LoginAttemptDB) GetLoginsByUserId(ctx context.Context, userID int, limit, offset int) ([]*models.LoginAttempt, error) {
	query := `
		SELECT id, user_id, ip_address, user_agent, successful, created_at
		FROM login_attempts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := db.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.LoginAttempt
	for rows.Next() {
		attempt := &models.LoginAttempt{}
		err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.IPAddress,
			&attempt.UserAgent,
			&attempt.Successful,
			&attempt.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attempts, nil
}
