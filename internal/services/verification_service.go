package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/sendgrid/rest"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sethvargo/go-password/password"
)

const PinLifeTime = 5

type EmailSender interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type VerificationService interface {
	SendVerificationEmail(ctx context.Context, userId int, email string) error
	GetVerification(ctx context.Context, id int) (*models.UserVerification, error)
	DeleteByUserId(ctx context.Context, userId int) error
	UpdatePin(ctx context.Context, userId int, email string) error
}

type verificationService struct {
	verificationRepo repo.VerificationRepository
	emailClient      EmailSender
}

func NewVerificationService(verificationRepo repo.VerificationRepository, emailClient EmailSender) *verificationService {
	return &verificationService{
		verificationRepo: verificationRepo,
		emailClient:      emailClient,
	}
}

func (s *verificationService) SendVerificationEmail(ctx context.Context, userId int, email string) error {
	pin, err := password.Generate(6, 2, 0, false, true)
	if err != nil {
		return err
	}
	println(userId)
	verification := &models.UserVerification{
		UserEmail:       email,
		UserId:          userId,
		VerificationPin: pin,
		PinExpiration:   time.Now().Add(PinLifeTime * time.Minute),
	}

	id, err := s.verificationRepo.AddPendingVerification(ctx, verification)
	if err != nil {
		return err
	}

	message := mail.NewV3MailInit(
		mail.NewEmail("ClassConnect service", "bmorseletto@fi.uba.ar"),
		"Verification Code",
		mail.NewEmail("User", email),
		mail.NewContent("text/plain", fmt.Sprintf("Your verification code is %d-%s", id, pin)),
	)

	if _, err := s.emailClient.Send(message); err != nil {
		return err
	}

	return nil
}

func (s *verificationService) GetVerification(ctx context.Context, id int) (*models.UserVerification, error) {
	return s.verificationRepo.GetVerificationById(ctx, id)
}

func (s *verificationService) GetVerificationByEmail(ctx context.Context, email string) (*models.UserVerification, error) {
	return s.verificationRepo.GetVerificationByEmail(ctx, email)
}

func (s *verificationService) DeleteByUserId(ctx context.Context, userId int) error {
	return s.verificationRepo.DeleteByUserId(ctx, userId)
}

func (s *verificationService) UpdatePin(ctx context.Context, userId int, email string) error {
	pin, err := password.Generate(6, 2, 0, false, true)
	if err != nil {
		return err
	}

	verification, err := s.GetVerificationByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return s.SendVerificationEmail(ctx, userId, email)
		}
		return err
	}

	err = s.verificationRepo.UpdatePin(ctx, verification.Id, pin)
	if err != nil {
		return err
	}

	message := mail.NewV3MailInit(
		mail.NewEmail("ClassConnect service", "bmorseletto@fi.uba.ar"),
		"Verification Code",
		mail.NewEmail("User", email),
		mail.NewContent("text/plain", fmt.Sprintf("Your verification code is %d-%s", verification.Id, pin)),
	)

	if _, err := s.emailClient.Send(message); err != nil {
		return err
	}

	return nil
}
