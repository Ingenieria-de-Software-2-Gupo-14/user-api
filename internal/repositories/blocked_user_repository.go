package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/errors"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

type BlockedUserRepository interface {
	BlockUser(ctx context.Context, blockedUserID int, reason string, blockerID *int, blockedUntil *time.Time) error
	UnblockUser(ctx context.Context, blockedUserID int) error
	GetBlocksByUserId(ctx context.Context, userID int) ([]models.BlockedUser, error)
}

type BlockedUserDB struct {
	DB *sql.DB
}

func NewBlockedUserRepository(db *sql.DB) *BlockedUserDB {
	return &BlockedUserDB{DB: db}
}

func (db *BlockedUserDB) BlockUser(ctx context.Context, blockedUserID int, reason string, blockerID *int, blockedUntil *time.Time) error {
	query := `
		INSERT INTO blocked_users (blocked_user_id, reason, blocker_id, blocked_until, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := db.DB.ExecContext(ctx, query, blockedUserID, reason, blockerID, blockedUntil, time.Now())
	return err
}

func (db *BlockedUserDB) UnblockUser(ctx context.Context, blockedUserID int) error {
	query := `
		UPDATE blocked_users
		SET blocked_until = NOW(), reason = reason || ' [UNBLOCK]'
		WHERE blocked_user_id = $1 AND (blocked_until IS NULL OR blocked_until > NOW())
		`
	_, err := db.DB.ExecContext(ctx, query, blockedUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound // No se encontró el bloqueo activo
		}
		return err // Otro error de base de datos
	}

	return nil // Desbloqueo exitoso
}

func (db *BlockedUserDB) GetBlocksByUserId(ctx context.Context, userID int) ([]models.BlockedUser, error) {
	query := `
		SELECT id, created_at, blocked_until, reason, blocker_id, blocked_user_id
		FROM blocked_users
		WHERE blocked_user_id = $1
		ORDER BY created_at DESC`

	rows, err := db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blocks []models.BlockedUser
	for rows.Next() {
		var block models.BlockedUser
		if err := rows.Scan(&block.Id, &block.CreatedAt, &block.BlockedUntil, &block.Reason, &block.BlockerId, &block.BlockedUserId); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}
