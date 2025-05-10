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
	AddPendingVerification(ctx context.Context, user *models.UserVerification) (int, error)
	GetPendingVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error)
	DeleteByEmail(ctx context.Context, email string) error
	UpdatePin(ctx context.Context, email string, pin string) error
	DeleteExpired() error
}

type verificationRepository struct {
	DB *sql.DB
}

func CreateVerificationRepo(db *sql.DB) *verificationRepository {
	return &verificationRepository{DB: db}
}

func (db verificationRepository) AddPendingVerification(ctx context.Context, user *models.UserVerification) (int, error) {
	query := `
		INSERT INTO verification (name, surname, password, email, verification_pin, pin_expiration)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id int
	err := db.DB.QueryRowContext(ctx, query,
		user.Name, user.Surname, user.Password, user.Email, user.VerificationPin, user.PinExpiration,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db verificationRepository) GetPendingVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error) {
	query := `SELECT id,email,name, surname, password,  verification_pin, pin_expiration FROM verification 
		WHERE email ILIKE $1`
	row := db.DB.QueryRowContext(ctx, query, email)
	var verification models.UserVerification

	err := row.Scan(
		&verification.Id, &verification.Email, &verification.Name, &verification.Surname, &verification.Password, &verification.VerificationPin, &verification.PinExpiration,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &verification, err
}

func (db verificationRepository) DeleteByEmail(ctx context.Context, email string) error {
	_, err := db.DB.ExecContext(ctx, "DELETE FROM verification WHERE email ILIKE $1", email)
	return err
}

func (db verificationRepository) UpdatePin(ctx context.Context, email string, pin string) error {
	rows, err := db.DB.ExecContext(ctx, "UPDATE verification SET verification_pin = $2, pin_expiration = $3  WHERE email ILIKE $1", email, pin, time.Now().Add(PinLifeTime*time.Minute))
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

func (db verificationRepository) DeleteExpired() error {
	rows, err := db.DB.Exec(`DELETE FROM verification WHERE pin_expiration < NOW()`)
	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	println(affected)
	return err
}
