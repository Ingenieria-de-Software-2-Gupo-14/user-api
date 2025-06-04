package services

import (
	"context"
	"time"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	repo "github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/repositories"
)

type RulesService interface {
	CreateRule(ctx context.Context, rule models.Rule, userId int) error
	DeleteRule(ctx context.Context, ruleId int, userId int) error
	GetRules(ctx context.Context) ([]models.Rule, error)
	ModifyRule(ctx context.Context, ruleId int, modification models.RuleModify, userId int) error
	GetAudits(ctx context.Context) ([]models.Audit, error)
}

type rulesService struct {
	rulesRepo repo.RulesRepository
}

// NewRulesService creates and returns a database
func NewRulesService(rulesRepo repo.RulesRepository) *rulesService {
	return &rulesService{rulesRepo}
}

func (s rulesService) CreateRule(ctx context.Context, rule models.Rule, userId int) error {
	if rule.EffectiveDate.IsZero() {
		rule.EffectiveDate = time.Now()
	}
	return s.rulesRepo.AddRule(ctx, rule, userId)
}

func (s rulesService) DeleteRule(ctx context.Context, ruleId int, userId int) error {
	return s.rulesRepo.DeleteRule(ctx, ruleId, userId)
}

func (s rulesService) GetRules(ctx context.Context) ([]models.Rule, error) {
	return s.rulesRepo.GetRules(ctx)
}
func (s rulesService) GetAudits(ctx context.Context) ([]models.Audit, error) {
	return s.rulesRepo.GetAudit(ctx)
}

func (s rulesService) ModifyRule(ctx context.Context, ruleId int, modification models.RuleModify, userId int) error {
	return s.rulesRepo.ModifyRule(ctx, ruleId, modification, userId)
}
