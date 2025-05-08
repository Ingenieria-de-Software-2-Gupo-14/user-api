package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

	mock.ExpectExec(`UPDATE blocked_users SET blocked_until = NOW\(\), reason = reason \|\| ' \[UNBLOCK\]' WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\)`).
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

	mock.ExpectExec(`UPDATE blocked_users SET blocked_until = NOW\(\), reason = reason \|\| ' \[UNBLOCK\]' WHERE blocked_user_id = \$1 AND \(blocked_until IS NULL OR blocked_until > NOW\(\)\)`).
		WithArgs(blockedUserID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	err = repo.UnblockUser(ctx, blockedUserID)
	assert.Equal(t, ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_GetBlocksByUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	now := time.Now()
	blockedUntil := now.Add(24 * time.Hour)

	rows := sqlmock.NewRows([]string{"id", "created_at", "blocked_until", "reason", "blocker_id", "blocked_user_id"}).
		AddRow(1, now, blockedUntil, "Violation of terms", 2, userID).
		AddRow(2, now.Add(-48*time.Hour), now.Add(-24*time.Hour), "Previous violation [UNBLOCK]", 3, userID)

	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	blocks, err := repo.GetBlocksByUserId(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, blocks, 2)
	assert.Equal(t, 1, blocks[0].Id)
	assert.Equal(t, userID, blocks[0].BlockedUserId)
	assert.Equal(t, "Violation of terms", blocks[0].Reason)
	assert.Equal(t, 2, *blocks[0].BlockerId)
	assert.Equal(t, 2, blocks[1].Id)
	assert.Equal(t, "Previous violation [UNBLOCK]", blocks[1].Reason)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_GetBlocksByUserId_NoBlocks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1

	rows := sqlmock.NewRows([]string{"id", "created_at", "blocked_until", "reason", "blocker_id", "blocked_user_id"})

	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnRows(rows)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	blocks, err := repo.GetBlocksByUserId(ctx, userID)
	assert.NoError(t, err)
	assert.Empty(t, blocks)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBlockedUserDB_GetBlocksByUserId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userID := 1
	expectedErr := sql.ErrConnDone

	mock.ExpectQuery(`SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id FROM blocked_users WHERE blocked_user_id = \$1 ORDER BY created_at DESC`).
		WithArgs(userID).
		WillReturnError(expectedErr)

	ctx := context.Background()
	repo := NewBlockedUserRepository(db)

	blocks, err := repo.GetBlocksByUserId(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, blocks)
	assert.NoError(t, mock.ExpectationsWereMet())
}
