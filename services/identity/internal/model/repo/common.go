package repo

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	identityUsersUsernameUniqueConstraint = "ux_identity_users_username"
	identityUsersEmailUniqueConstraint    = "ux_identity_users_email"
)

// UniqueViolationKind is the semantic kind of a PostgreSQL unique violation.
// UniqueViolationKind 表示 PostgreSQL 唯一冲突的语义类型。
type UniqueViolationKind string

const (
	UniqueViolationUnknown  UniqueViolationKind = "unknown_conflict"
	UniqueViolationUsername UniqueViolationKind = "username_conflict"
	UniqueViolationEmail    UniqueViolationKind = "email_conflict"
)

// IsNotFound reports whether the error means the record does not exist.
// IsNotFound 判断错误是否表示记录不存在。
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsUniqueViolation reports whether the error is a PostgreSQL unique violation.
// IsUniqueViolation 判断错误是否为 PostgreSQL 唯一约束冲突。
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

// UniqueViolationConstraint returns the PostgreSQL unique constraint name when available.
// UniqueViolationConstraint 返回 PostgreSQL 唯一约束名（若可获取）。
func UniqueViolationConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.ConstraintName
	}

	return ""
}

// ParseUniqueViolation classifies PostgreSQL unique violations into semantic conflict kinds.
// ParseUniqueViolation 将 PostgreSQL 唯一冲突解析为语义化冲突类型。
func ParseUniqueViolation(err error) (UniqueViolationKind, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return UniqueViolationUnknown, false
	}

	switch strings.ToLower(strings.TrimSpace(pgErr.ConstraintName)) {
	case identityUsersUsernameUniqueConstraint:
		return UniqueViolationUsername, true
	case identityUsersEmailUniqueConstraint:
		return UniqueViolationEmail, true
	}

	detail := strings.ToLower(pgErr.Detail)
	message := strings.ToLower(pgErr.Message)
	combined := detail + " " + message
	if strings.Contains(combined, "(username)") || strings.Contains(combined, "username") {
		return UniqueViolationUsername, true
	}
	if strings.Contains(combined, "(email)") || strings.Contains(combined, "email") {
		return UniqueViolationEmail, true
	}

	return UniqueViolationUnknown, true
}
