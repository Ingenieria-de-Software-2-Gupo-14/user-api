package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewBlockedUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewBlockedUserRepository(db)
	assert.NotNil(t, repo)
}

func TestBlockedUserDB_BlockUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockedUserID := 1
	reason := "Violation of terms"
	blockerID := 2
	now := time.Now()
	blockedUntil := now.Add(24 * time.Hour)

	mock.ExpectExec(`INSERT INTO blocked_users \(blocked_user_id, reason, blocker_id, blocked_until, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(blockedUserID, reason, &blockerID, &blockedUntil, sqlmock.AnyArg()). // Using AnyArg for created_at since it's generated at runtime
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	err = repo.BlockUser(ctx, blockedUserID, reason, &blockerID, &blockedUntil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_BlockUser_PermanentBlock(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockedUserID := 1
	reason := "Permanent block"
	var blockerID *int = nil
	var blockedUntil *time.Time = nil

	mock.ExpectExec(`INSERT INTO blocked_users \(blocked_user_id, reason, blocker_id, blocked_until, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(blockedUserID, reason, blockerID, blockedUntil, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	err = repo.BlockUser(ctx, blockedUserID, reason, blockerID, blockedUntil)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_UnblockUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockedUserID := 1

	mock.ExpectExec(`UPDATE blocked_users SET blocked_until = NOW\(\) WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\)`).
		WithArgs(blockedUserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	err = repo.UnblockUser(ctx, blockedUserID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_UnblockUser_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	blockedUserID := 1

	mock.ExpectExec(`UPDATE blocked_users SET blocked_until = NOW\(\) WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\)`).
		WithArgs(blockedUserID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	err = repo.UnblockUser(ctx, blockedUserID)
	assert.Equal(t, errors.ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_IsUserBlocked_Blocked(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	blockID := 5
	createdAt := time.Now().Add(-24 * time.Hour)
	blockedUntil := time.Now().Add(24 * time.Hour)
	reason := "Terms violation"
	blockerID := 2

	columns := []string{"id", "created_at", "blocked_until", "reason", "blocker_id", "blocked_user_id"}
	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\) ORDER BY created_at DESC LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(
			blockID, createdAt, &blockedUntil, reason, &blockerID, userID))

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	isBlocked, blockedInfo, err := repo.IsUserBlocked(ctx, userID)
	assert.NoError(t, err)
	assert.True(t, isBlocked)
	assert.NotNil(t, blockedInfo)
	assert.Equal(t, blockID, blockedInfo.Id)
	assert.Equal(t, createdAt.Unix(), blockedInfo.CreatedAt.Unix()) // Compare Unix timestamps to avoid time precision issues
	assert.NotNil(t, blockedInfo.BlockedUntil)
	assert.Equal(t, blockedUntil.Unix(), blockedInfo.BlockedUntil.Unix())
	assert.Equal(t, reason, blockedInfo.Reason)
	assert.NotNil(t, blockedInfo.BlockerId)
	assert.Equal(t, blockerID, *blockedInfo.BlockerId)
	assert.Equal(t, userID, blockedInfo.BlockedUserId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_IsUserBlocked_NotBlocked(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1

	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\) ORDER BY created_at DESC LIMIT 1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	isBlocked, blockedInfo, err := repo.IsUserBlocked(ctx, userID)
	assert.NoError(t, err)
	assert.False(t, isBlocked)
	assert.Nil(t, blockedInfo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_IsUserBlocked_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	expectedErr := sql.ErrConnDone

	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\) ORDER BY created_at DESC LIMIT 1`).
		WithArgs(userID).
		WillReturnError(expectedErr)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	isBlocked, blockedInfo, err := repo.IsUserBlocked(ctx, userID)
	assert.Equal(t, expectedErr, err)
	assert.False(t, isBlocked)
	assert.Nil(t, blockedInfo)
	assert.NoError(t, mock.ExpectationsWereMet())
}
