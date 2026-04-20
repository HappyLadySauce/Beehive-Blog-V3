package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestCredentialLocalRepository verifies local credential persistence.
// TestCredentialLocalRepository 验证本地凭证持久化。
func TestCredentialLocalRepository(t *testing.T) {
	t.Parallel()

	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	ctx := context.Background()

	credential := &entity.CredentialLocal{
		UserID:            user.ID,
		PasswordHash:      "hash",
		PasswordUpdatedAt: time.Now().UTC(),
	}
	if err := store.CredentialLocals.Create(ctx, credential); err != nil {
		t.Fatalf("expected create credential to succeed, got %v", err)
	}

	got, err := store.CredentialLocals.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected get credential to succeed, got %v", err)
	}
	if got.UserID != user.ID || got.PasswordHash != "hash" {
		t.Fatalf("unexpected credential: %+v", got)
	}

	if _, err := store.CredentialLocals.GetByUserID(ctx, user.ID+999); err == nil {
		t.Fatalf("expected missing credential lookup to fail")
	}

	_ = auth.UserRoleMember
}
