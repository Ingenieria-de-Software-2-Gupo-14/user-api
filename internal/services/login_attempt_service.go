package services

import (
	"context"
	"ing-soft-2-tp1/internal/models"
	"ing-soft-2-tp1/internal/repositories"
	"time"
)

const (
	// MaxFailedAttempts is the maximum number of failed login attempts before blocking the user
	MaxFailedAttempts = 5
	// BlockDuration is the duration in minutes for which the user will be blocked
	BlockDuration = 10
	// BlockedUserReason is the reason for blocking the user
	BlockedUserReason = "Too many failed login attempts"
)

type LoginAttemptService interface {
	AddLoginAttempt(ctx context.Context, userID int, ipAddress, userAgent string, successful bool) error
	GetLoginsByUserId(ctx context.Context, userID int, limit, offset int) ([]*models.LoginAttempt, error)
}

type LoginAttemptServiceImpl struct {
	repo            repositories.LoginAttemptRepository
	blockedUserRepo repositories.BlockedUserRepository
}

func NewLoginAttemptService(
	repo repositories.LoginAttemptRepository,
	blockedUserRepo repositories.BlockedUserRepository,
) LoginAttemptService {

	return &LoginAttemptServiceImpl{
		repo:            repo,
		blockedUserRepo: blockedUserRepo,
	}
}

// AddLoginAttempt logs a login attempt and checks if the user should be blocked
// based on the number of failed attempts in the last 10 minutes.
func (s *LoginAttemptServiceImpl) AddLoginAttempt(ctx context.Context, userID int, ipAddress, userAgent string, successful bool) error {
	// Registrar el intento de inicio de sesión
	err := s.repo.AddLoginAttempt(ctx, userID, ipAddress, userAgent, successful)
	if err != nil {
		return err
	}

	if successful {
		return nil
	}

	// Verificar si el usuario debe ser bloqueado
	attempts, err := s.repo.BadLoginAttemptsInLast10Minutes(ctx, userID)
	if err != nil {
		return err
	}

	if attempts >= MaxFailedAttempts {
		blockedUntil := time.Now().Add(BlockDuration * time.Minute)
		return s.blockedUserRepo.BlockUser(ctx, userID, BlockedUserReason, nil, &blockedUntil)
	}

	return nil
}

// GetLoginsByUserId obtiene el historial de intentos de login para un usuario específico
func (s *LoginAttemptServiceImpl) GetLoginsByUserId(ctx context.Context, userID int, limit, offset int) ([]*models.LoginAttempt, error) {
	return s.repo.GetLoginsByUserId(ctx, userID, limit, offset)
}
