package repo

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	ConstraintItemSlug = "ux_content_items_slug"
	ConstraintTagName  = "ux_content_tags_name"
	ConstraintTagSlug  = "ux_content_tags_slug"
)

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func UniqueConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return strings.TrimSpace(pgErr.ConstraintName)
	}
	return ""
}
