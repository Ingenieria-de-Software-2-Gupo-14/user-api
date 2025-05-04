package services

import (
	"context"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/repositories"
	"time"
)

type BlockedUserService interface {
	BlockUser(ctx context.Context, blockedUserID int, reason string, blockerID *int, blockedUntil *time.Time) error
	UnblockUser(ctx context.Context, blockedUserID int) error
	IsUserBlocked(ctx context.Context, userID int) (bool, *models.BlockedUser, error)
}

type BlockedUserServiceImpl struct {
	repo repositories.BlockedUserRepository
}

func NewBlockedUserService(repo repositories.BlockedUserRepository) BlockedUserService {
	return &BlockedUserServiceImpl{
		repo: repo,
	}
}

// BlockUser bloquea un usuario por la razón especificada
func (s *BlockedUserServiceImpl) BlockUser(ctx context.Context, blockedUserID int, reason string, blockerID *int, blockedUntil *time.Time) error {
	return s.repo.BlockUser(ctx, blockedUserID, reason, blockerID, blockedUntil)
}

// UnblockUser desbloquea un usuario
func (s *BlockedUserServiceImpl) UnblockUser(ctx context.Context, blockedUserID int) error {
	return s.repo.UnblockUser(ctx, blockedUserID)
}

// IsUserBlocked verifica si un usuario está bloqueado y devuelve la información de bloqueo
func (s *BlockedUserServiceImpl) IsUserBlocked(ctx context.Context, userID int) (bool, *models.BlockedUser, error) {
	return s.repo.IsUserBlocked(ctx, userID)
}
