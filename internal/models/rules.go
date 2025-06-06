package models

import (
	"database/sql"
	"time"
)

type Rule struct {
	Id                   int       `json:"id"`
	Title                string    `json:"Title" binding:"required"`
	Description          string    `json:"Description" binding:"required"`
	EffectiveDate        time.Time `json:"effectiveDate" `
	ApplicationCondition string    `json:"ApplicationCondition" binding:"required"`
}

type RuleModify struct {
	Title                string `json:"Title"`
	Description          string `json:"Description" `
	ApplicationCondition string `json:"ApplicationCondition" `
}

type Audit struct {
	Id                   int
	RuleId               sql.NullInt64
	UserId               sql.NullInt64
	ModificationDate     time.Time
	NatureOfModification string
}

type AuditData struct {
	Id                   int
	RuleId               int
	UserId               int
	ModificationDate     time.Time
	NatureOfModification string
}
