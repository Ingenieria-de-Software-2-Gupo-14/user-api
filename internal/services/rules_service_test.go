package services

import (
	"context"
	"database/sql"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNewRulesService(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	assert.NotNil(t, service)
}

func TestRulesService_CreateRule(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	c := context.Background()

	userId := 1
	rule := models.Rule{
		Id:                   userId,
		Title:                "title",
		Description:          "description",
		EffectiveDate:        time.Time{},
		ApplicationCondition: "condition",
	}

	mockRepo.EXPECT().AddRule(c, mock.Anything, userId).Return(nil)

	err := service.CreateRule(c, rule, userId)
	assert.NoError(t, err)
}

func TestRulesService_DeleteRule(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	c := context.Background()

	userId := 1
	ruleId := 1
	mockRepo.EXPECT().DeleteRule(c, ruleId, userId).Return(nil)

	err := service.DeleteRule(c, ruleId, userId)
	assert.NoError(t, err)
}

func TestRulesService_GetRules(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	c := context.Background()

	rule := models.Rule{
		Id:                   1,
		Title:                "title",
		Description:          "desc",
		EffectiveDate:        time.Time{},
		ApplicationCondition: "condition",
	}

	rules := []models.Rule{rule}

	mockRepo.EXPECT().GetRules(c).Return(rules, nil)

	rulesResult, err := service.GetRules(c)
	assert.NoError(t, err)
	assert.Equal(t, rules, rulesResult)
}

func TestRulesService_GetAudits(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	c := context.Background()

	audit := models.Audit{
		Id:                   1,
		RuleId:               sql.NullInt64{Int64: 1, Valid: true},
		UserId:               sql.NullInt64{Int64: 1, Valid: true},
		ModificationDate:     time.Time{},
		NatureOfModification: "modification",
	}

	audits := []models.Audit{audit}
	mockRepo.EXPECT().GetAudit(c).Return(audits, nil)

	auditsResult, err := service.GetAudits(c)
	assert.NoError(t, err)
	assert.Equal(t, audits, auditsResult)
}

func TestRulesService_ModifyRule(t *testing.T) {
	mockRepo := repositories.NewMockRulesRepository(t)
	service := NewRulesService(mockRepo)
	c := context.Background()

	userId := 1
	ruleId := 1
	modification := models.RuleModify{
		Title:                "title",
		Description:          "description",
		ApplicationCondition: "Application",
	}

	mockRepo.EXPECT().ModifyRule(c, ruleId, modification, userId).Return(nil)

	err := service.ModifyRule(c, ruleId, modification, userId)
	assert.NoError(t, err)
}
