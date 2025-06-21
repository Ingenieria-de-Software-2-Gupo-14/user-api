package repositories

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateVerificationRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	result := CreateVerificationRepo(db)

	assert.NotNil(t, result)
}

func TestVerificationRepository_AddPendingVerification(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)

	ctx := context.Background()
	expectedID := 1
	verification := &models.UserVerification{
		UserId:          1,
		UserEmail:       "user@example.com",
		VerificationPin: "123456",
		PinExpiration:   time.Now().Add(10 * time.Minute),
	}

	mock.ExpectQuery(`INSERT INTO verification \(user_id, user_email, verification_pin, pin_expiration\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
		WithArgs(
			verification.UserId,
			verification.UserEmail,
			verification.VerificationPin,
			verification.PinExpiration,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := repo.AddPendingVerification(ctx, verification)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_DeleteByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)

	ctx := context.Background()
	userId := 1

	mock.ExpectExec(`DELETE FROM verification WHERE user_id = \$1`).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(0, 1)) // simulate 1 row deleted

	err = repo.DeleteByUserId(ctx, userId)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_GetVerificationByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	email := "test@test.com"
	now := time.Now()

	expected := &models.UserVerification{
		Id:              1,
		UserId:          2,
		UserEmail:       email,
		VerificationPin: "123456",
		PinExpiration:   now.Add(10 * time.Minute),
		CreatedAt:       now,
	}

	mock.ExpectQuery(`SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE user_email ILIKE \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "user_email", "verification_pin", "pin_expiration", "created_at",
		}).AddRow(
			expected.Id,
			expected.UserId,
			expected.UserEmail,
			expected.VerificationPin,
			expected.PinExpiration,
			expected.CreatedAt,
		))

	result, err := repo.GetVerificationByEmail(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_GetVerificationByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	email := "test@test.com"

	mock.ExpectQuery(`SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE user_email ILIKE \$1`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetVerificationByEmail(ctx, email)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_GetVerificationById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	now := time.Now()

	email := "test@test.com"

	expected := &models.UserVerification{
		Id:              1,
		UserId:          2,
		UserEmail:       email,
		VerificationPin: "123456",
		PinExpiration:   now.Add(10 * time.Minute),
		CreatedAt:       now,
	}

	mock.ExpectQuery(`SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE id = \$1`).
		WithArgs(expected.Id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "user_email", "verification_pin", "pin_expiration", "created_at",
		}).AddRow(
			expected.Id, expected.UserId, expected.UserEmail, expected.VerificationPin,
			expected.PinExpiration, expected.CreatedAt,
		))

	result, err := repo.GetVerificationById(ctx, expected.Id)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_GetVerificationById_NotFounf(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	id := 99

	mock.ExpectQuery(`SELECT id, user_id, user_email, verification_pin, pin_expiration, created_at FROM verification
		WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetVerificationById(ctx, id)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_UpdatePin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	id := 1
	pin := "654321"

	mock.ExpectExec(`UPDATE verification SET verification_pin = \$2, pin_expiration = \$3 WHERE id = \$1`).
		WithArgs(id, pin, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.UpdatePin(ctx, id, pin)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestVerificationRepository_UpdatePin_NoRowAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateVerificationRepo(db)
	ctx := context.Background()
	id := 1
	pin := "654321"

	mock.ExpectExec(`UPDATE verification SET verification_pin = \$2, pin_expiration = \$3 WHERE id = \$1`).
		WithArgs(id, pin, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.UpdatePin(ctx, id, pin)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
