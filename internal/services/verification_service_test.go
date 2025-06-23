package services

import (
	"context"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewVerificationService(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	mockEmail := NewMockEmailSender(t)
	service := NewVerificationService(mockRepo, mockEmail)

	assert.NotNil(t, service)
}

func TestVerificationService_GetVerification(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	service := NewVerificationService(mockRepo, nil)

	ctx := context.Background()
	userID := 1
	expected := &models.UserVerification{
		Id:              1,
		UserId:          userID,
		UserEmail:       "test@email.com",
		VerificationPin: "123456",
		PinExpiration:   time.Time{},
		CreatedAt:       time.Time{},
	}

	mockRepo.EXPECT().
		GetVerificationById(ctx, userID).
		Return(expected, nil)

	result, err := service.GetVerification(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestVerificationService_GetVerificationByEmail(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	service := NewVerificationService(mockRepo, nil)

	ctx := context.Background()
	email := "test@email.com"
	userID := 1
	expected := &models.UserVerification{
		Id:              1,
		UserId:          userID,
		UserEmail:       email,
		VerificationPin: "123456",
		PinExpiration:   time.Time{},
		CreatedAt:       time.Time{},
	}

	mockRepo.EXPECT().
		GetVerificationByEmail(ctx, email).
		Return(expected, nil)

	result, err := service.GetVerificationByEmail(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestVerificationService_DeleteByUserId(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	service := NewVerificationService(mockRepo, nil)

	ctx := context.Background()
	userID := 1

	mockRepo.EXPECT().
		DeleteByUserId(ctx, userID).
		Return(nil)

	err := service.DeleteByUserId(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSendVerificationEmail(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	mockEmail := NewMockEmailSender(t)

	service := NewVerificationService(mockRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	email := "user@test.com"

	mockRepo.EXPECT().
		AddPendingVerification(ctx, mock.AnythingOfType("*models.UserVerification")).
		Return(123, nil)

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	err := service.SendVerificationEmail(ctx, userID, email)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestUpdatePin(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	mockEmail := NewMockEmailSender(t)
	service := NewVerificationService(mockRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	email := "user@test.com"
	verification := &models.UserVerification{
		Id:        42,
		UserId:    userID,
		UserEmail: email,
	}

	mockRepo.EXPECT().
		GetVerificationByEmail(ctx, email).
		Return(verification, nil)

	mockRepo.EXPECT().
		UpdatePin(ctx, verification.Id, mock.AnythingOfType("string")).
		Return(nil)

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	err := service.UpdatePin(ctx, userID, email)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestUpdatePin_VerificationEmailNotFound(t *testing.T) {
	mockRepo := repo.NewMockVerificationRepository(t)
	mockEmail := NewMockEmailSender(t)
	service := NewVerificationService(mockRepo, mockEmail)

	ctx := context.Background()
	userID := 1
	email := "user@test.com"

	mockRepo.EXPECT().
		GetVerificationByEmail(ctx, email).
		Return(nil, repo.ErrNotFound)

	mockRepo.EXPECT().
		AddPendingVerification(ctx, mock.AnythingOfType("*models.UserVerification")).
		Return(123, nil)

	mockEmail.On("Send", mock.AnythingOfType("*mail.SGMailV3")).
		Return(&rest.Response{StatusCode: 202}, nil)

	err := service.UpdatePin(ctx, userID, email)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}
