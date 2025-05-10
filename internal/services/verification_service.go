package services

import (
	"context"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/sethvargo/go-password/password"
)

const PinLifeTime = 5

type VerificationService interface {
	CreatePendingVerification(ctx context.Context, request models.CreateUserRequest, admin bool) (string, error)
	GetPendingVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error)
	DeleteByEmail(ctx context.Context, email string) error
}

type verificationService struct {
	verificationRepo repo.VerificationRepository
}

func NewVerificationService(verificationRepo repo.VerificationRepository) *verificationService {
	return &verificationService{verificationRepo}
}

func (s *verificationService) CreatePendingVerification(ctx context.Context, request models.CreateUserRequest, admin bool) (string, error) {
	/*hashPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return "", err
	}*/
	pin, errPin := password.Generate(6, 2, 0, false, true)
	if errPin != nil {
		return "", errPin
	}
	user := &models.UserVerification{
		Email:           request.Email,
		Name:            request.Name,
		Surname:         request.Surname,
		Password:        request.Password,
		VerificationPin: pin,
		PinExpiration:   time.Now().Add(PinLifeTime * time.Minute),
	}
	_, err := s.verificationRepo.AddPendingVerification(ctx, user)
	if err != nil {
		return "", err
	}
	return pin, err
}

func (s *verificationService) GetPendingVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error) {
	return s.verificationRepo.GetPendingVerificationByEmail(ctx, email)
}

func (s *verificationService) DeleteByEmail(ctx context.Context, email string) error {
	return s.verificationRepo.DeleteByEmail(ctx, email)
}
