package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Ingenieria-de-Software-2-Gupo-14/user-api/internal/models"

	_ "github.com/lib/pq"
)

type RulesRepository interface {
	AddRule(ctx context.Context, rule models.Rule, userId int) error
	DeleteRule(ctx context.Context, ruleId int, userId int) error
	GetRules(ctx context.Context) ([]models.Rule, error)
	ModifyRule(ctx context.Context, ruleId int, modification models.RuleModify, userId int) error
	GetAudit(ctx context.Context) ([]models.Audit, error)
}

type rulesRepository struct {
	DB *sql.DB
}

// CreateRulesRepo creates and returns a database
func CreateRulesRepo(db *sql.DB) *rulesRepository {
	return &rulesRepository{DB: db}
}

func (db rulesRepository) AddRule(ctx context.Context, rule models.Rule, userId int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	query := `
		INSERT INTO rules (title, description,effective_date, application_condition)
		VALUES ($1, $2, $3,$4)
		RETURNING id`
	var id int
	err = tx.QueryRowContext(ctx, query,
		&rule.Title, &rule.Description, &rule.EffectiveDate, &rule.ApplicationCondition,
	).Scan(&id)

	if err != nil {
		return err
	}
	_, err = tx.QueryContext(ctx, `
	 		INSERT INTO rules_audit (rule_id, user_id, nature_of_modification)
	 		VALUES ($1, $2, $3)`, id, userId, fmt.Sprint("Created Rule ", id, ": ", rule.Title))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db rulesRepository) DeleteRule(ctx context.Context, ruleId int, userId int) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var deletedRule models.Rule
	err = tx.QueryRow(`
			DELETE FROM rules
			WHERE id = $1
			RETURNING title, description, application_condition
		`, ruleId).Scan(&deletedRule.Title, &deletedRule.Description, &deletedRule.ApplicationCondition)
	if err != nil {
		return err
	}
	_, err = tx.QueryContext(ctx, `
	 		INSERT INTO rules_audit (rule_id, user_id, nature_of_modification)
	 		VALUES ($1, $2, $3)`, nil, userId, fmt.Sprint("Deleted Rule ", ruleId, ": ", deletedRule.Title))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db rulesRepository) GetRules(ctx context.Context) ([]models.Rule, error) {
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, title, description, effective_date, application_condition
		FROM rules`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []models.Rule
	for rows.Next() {
		var rule models.Rule
		err := rows.Scan(&rule.Id, &rule.Title, &rule.Description, &rule.EffectiveDate, &rule.ApplicationCondition)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (db rulesRepository) ModifyRule(ctx context.Context, ruleId int, modification models.RuleModify, userId int) error {
	modification_nature := fmt.Sprint("Modified Rule ", ruleId, ":")
	counter := 1
	query := "UPDATE rules SET"
	params := []interface{}{}
	setClauses := []string{}
	if modification.Title != "" {
		setClauses = append(setClauses, fmt.Sprint("title = $", counter))
		counter += 1
		params = append(params, modification.Title)
		modification_nature += " Changed title to: " + modification.Title
	}
	if modification.Description != "" {
		setClauses = append(setClauses, fmt.Sprint("description = $", counter))
		counter += 1
		params = append(params, modification.Description)
		modification_nature += " Changed description to: " + modification.Description
	}
	if modification.ApplicationCondition != "" {
		setClauses = append(setClauses, fmt.Sprint("application_condition = $", counter))
		counter += 1
		params = append(params, modification.ApplicationCondition)
		modification_nature += " Changed application condition to: " + modification.ApplicationCondition
	}
	if len(setClauses) == 0 {
		return nil
	}
	query += " " + strings.Join(setClauses, ", ") + fmt.Sprint(" WHERE id = $", counter)
	params = append(params, ruleId)
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}
	_, err = tx.QueryContext(ctx, `
	 		INSERT INTO rules_audit (rule_id, user_id, nature_of_modification)
	 		VALUES ($1, $2, $3)`, ruleId, userId, modification_nature)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db rulesRepository) GetAudit(ctx context.Context) ([]models.Audit, error) {
	rows, err := db.DB.QueryContext(ctx, `
		SELECT id, rule_id, user_id, modification_date, nature_of_modification
		FROM rules_audit`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var audits []models.Audit
	for rows.Next() {
		var audit models.Audit
		err := rows.Scan(&audit.Id, &audit.RuleId, &audit.UserId, &audit.ModificationDate, &audit.NatureOfModification)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}
	return audits, nil
}
