package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

const PinLifeTime = 5

type VerificationRepository interface {
	AddPendingVerification(ctx context.Context, verification *models.UserVerification) (int, error)
	GetVerificationById(ctx context.Context, id int) (*models.UserVerification, error)
	GetVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error)
	DeleteByUserId(ctx context.Context, userId int) error
	UpdatePin(ctx context.Context, userId int, pin string) error
}

type verificationRepository struct {
	DB *sql.DB
}

func CreateVerificationRepo(db *sql.DB) *verificationRepository {
	return &verificationRepository{DB: db}
}

func (db verificationRepository) AddPendingVerification(ctx context.Context, verification *models.UserVerification) (int, error) {
	query := `
		INSERT INTO verification (user_id, user_email, verification_pin, pin_expiration)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, query,
		verification.UserId, verification.UserEmail, verification.VerificationPin, verification.PinExpiration,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db verificationRepository) GetVerificationById(ctx context.Context, id int) (*models.UserVerification, error) {
	query := `SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE id = $1`
	row := db.DB.QueryRowContext(ctx, query, id)
	var verification models.UserVerification
	err := row.Scan(
		&verification.Id, &verification.UserId, &verification.UserEmail, &verification.VerificationPin,
		&verification.PinExpiration, &verification.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &verification, nil
}

func (db verificationRepository) GetVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error) {
	query := `SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE user_email ILIKE $1`
	row := db.DB.QueryRowContext(ctx, query, email)
	var verification models.UserVerification
	err := row.Scan(
		&verification.Id, &verification.UserId, &verification.UserEmail, &verification.VerificationPin,
		&verification.PinExpiration, &verification.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &verification, nil
}

func (db verificationRepository) DeleteByUserId(ctx context.Context, userId int) error {
	_, err := db.DB.ExecContext(ctx, "DELETE FROM verification WHERE user_id = $1", userId)
	return err
}

func (db verificationRepository) UpdatePin(ctx context.Context, userId int, pin string) error {
	rows, err := db.DB.ExecContext(ctx, "UPDATE verification SET verification_pin = $2, pin_expiration = $3 WHERE user_id = $1",
		userId, pin, time.Now().Add(PinLifeTime*time.Minute))
	if err != nil {
		return err
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return ErrNotFound
	}
	return err
}
