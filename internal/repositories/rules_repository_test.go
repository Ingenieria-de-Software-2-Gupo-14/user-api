package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateRulesRepo(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)
	assert.NotNil(t, repo)
}

func TestRulesRepository_AddRule(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)

	ctx := context.Background()
	rule := models.Rule{
		Title:                "Test Rule",
		Description:          "A rule for testing",
		EffectiveDate:        time.Now(),
		ApplicationCondition: "Condition A",
	}
	userId := 1

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO rules \(title, description,effective_date, application_condition\) VALUES \(\$1, \$2, \$3,\$4\) RETURNING id`).
		WithArgs(rule.Title, rule.Description, rule.EffectiveDate, rule.ApplicationCondition).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	mock.ExpectQuery(`INSERT INTO rules_audit \(rule_id, user_id, nature_of_modification\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(123, userId, fmt.Sprint("Created Rule ", 123, ": ", rule.Title)).
		WillReturnRows(sqlmock.NewRows([]string{"dummy"}))

	mock.ExpectCommit()

	err = repo.AddRule(ctx, rule, userId)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRulesRepository_DeleteRule(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)

	ctx := context.Background()
	ruleID := 1
	userID := 1
	deletedTitle := "Old Rule"
	deletedDescription := "No longer needed"
	deletedCondition := "Condition X"

	mock.ExpectBegin()

	mock.ExpectQuery(`DELETE FROM rules WHERE id = \$1 RETURNING title, description, application_condition`).
		WithArgs(ruleID).
		WillReturnRows(sqlmock.NewRows([]string{"title", "description", "application_condition"}).
			AddRow(deletedTitle, deletedDescription, deletedCondition))

	mock.ExpectQuery(`INSERT INTO rules_audit \(rule_id, user_id, nature_of_modification\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(nil, userID, fmt.Sprint("Deleted Rule ", ruleID, ": ", deletedTitle)).
		WillReturnRows(sqlmock.NewRows([]string{"dummy"}))

	mock.ExpectCommit()

	err = repo.DeleteRule(ctx, ruleID, userID)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRulesRepository_GetAudit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)

	ctx := context.Background()

	now := time.Now()
	expectedAudits := []models.Audit{
		{
			Id:                   1,
			RuleId:               sql.NullInt64{Int64: 10, Valid: true},
			UserId:               sql.NullInt64{Int64: 1, Valid: true},
			ModificationDate:     now,
			NatureOfModification: "Created Rule 10: Test Rule",
		},
		{
			Id:                   2,
			RuleId:               sql.NullInt64{Int64: 11, Valid: true},
			UserId:               sql.NullInt64{Int64: 1, Valid: true},
			ModificationDate:     now,
			NatureOfModification: "Deleted Rule 11: Old Rule",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "rule_id", "user_id", "modification_date", "nature_of_modification"}).
		AddRow(expectedAudits[0].Id, expectedAudits[0].RuleId, expectedAudits[0].UserId, expectedAudits[0].ModificationDate, expectedAudits[0].NatureOfModification).
		AddRow(expectedAudits[1].Id, expectedAudits[1].RuleId, expectedAudits[1].UserId, expectedAudits[1].ModificationDate, expectedAudits[1].NatureOfModification)

	mock.ExpectQuery(`SELECT id, rule_id, user_id, modification_date, nature_of_modification FROM rules_audit`).
		WillReturnRows(rows)

	audits, err := repo.GetAudit(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedAudits, audits)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRulesRepository_GetRules(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)

	ctx := context.Background()

	now := time.Now()
	expectedRules := []models.Rule{
		{
			Id:                   1,
			Title:                "Rule A",
			Description:          "First rule",
			EffectiveDate:        now,
			ApplicationCondition: "If condition A",
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "effective_date", "application_condition",
	}).
		AddRow(expectedRules[0].Id, expectedRules[0].Title, expectedRules[0].Description, expectedRules[0].EffectiveDate, expectedRules[0].ApplicationCondition)

	mock.ExpectQuery(`SELECT id, title, description, effective_date, application_condition FROM rules`).
		WillReturnRows(rows)

	rules, err := repo.GetRules(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedRules, rules)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRulesRepository_ModifyRule(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := CreateRulesRepo(db)

	ctx := context.Background()
	ruleID := 1
	userID := 1

	modification := models.RuleModify{
		Title:                "Updated Title",
		Description:          "Updated Description",
		ApplicationCondition: "Updated Condition",
	}

	modificationNature := fmt.Sprint("Modified Rule ", ruleID, ":")
	modificationNature += " Changed title to: " + modification.Title
	modificationNature += " Changed description to: " + modification.Description
	modificationNature += " Changed application condition to: " + modification.ApplicationCondition

	// Expect transaction begin
	mock.ExpectBegin()

	// Expect dynamic UPDATE
	mock.ExpectExec(`UPDATE rules SET title = \$1, description = \$2, application_condition = \$3 WHERE id = \$4`).
		WithArgs(modification.Title, modification.Description, modification.ApplicationCondition, ruleID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // dummy result

	// Expect audit insert
	mock.ExpectQuery(`INSERT INTO rules_audit \(rule_id, user_id, nature_of_modification\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(ruleID, userID, modificationNature).
		WillReturnRows(sqlmock.NewRows([]string{"dummy"}))

	// Expect commit
	mock.ExpectCommit()

	err = repo.ModifyRule(ctx, ruleID, modification, userID)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
