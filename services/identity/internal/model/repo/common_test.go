package repo

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// TestParseUniqueViolation verifies semantic classification for unique violations.
// TestParseUniqueViolation 验证唯一冲突的语义分类行为。
func TestParseUniqueViolation(t *testing.T) {
	t.Parallel()

	t.Run("username constraint", func(t *testing.T) {
		t.Parallel()

		kind, ok := ParseUniqueViolation(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: "ux_identity_users_username",
		})
		if !ok || kind != UniqueViolationUsername {
			t.Fatalf("expected username conflict, got kind=%s ok=%v", kind, ok)
		}
	})

	t.Run("email constraint", func(t *testing.T) {
		t.Parallel()

		kind, ok := ParseUniqueViolation(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: "ux_identity_users_email",
		})
		if !ok || kind != UniqueViolationEmail {
			t.Fatalf("expected email conflict, got kind=%s ok=%v", kind, ok)
		}
	})

	t.Run("fallback by detail", func(t *testing.T) {
		t.Parallel()

		kind, ok := ParseUniqueViolation(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: "renamed_constraint",
			Detail:         "Key (username)=(alice) already exists.",
		})
		if !ok || kind != UniqueViolationUsername {
			t.Fatalf("expected username conflict by detail, got kind=%s ok=%v", kind, ok)
		}
	})

	t.Run("unknown unique conflict", func(t *testing.T) {
		t.Parallel()

		kind, ok := ParseUniqueViolation(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: "ux_identity_unknown",
		})
		if !ok || kind != UniqueViolationUnknown {
			t.Fatalf("expected unknown unique conflict, got kind=%s ok=%v", kind, ok)
		}
	})

	t.Run("wrapped non unique error", func(t *testing.T) {
		t.Parallel()

		kind, ok := ParseUniqueViolation(errors.New("plain error"))
		if ok || kind != UniqueViolationUnknown {
			t.Fatalf("expected non unique result, got kind=%s ok=%v", kind, ok)
		}
	})
}
