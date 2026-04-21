package repo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
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
