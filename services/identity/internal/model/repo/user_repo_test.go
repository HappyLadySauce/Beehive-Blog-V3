package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestUserRepository verifies user persistence behavior.
// TestUserRepository 验证用户持久化行为。
func TestUserRepository(t *testing.T) {
	t.Parallel()

	store := testkit.NewStore(t)
	ctx := context.Background()
	email := "alice@example.com"
	user := &entity.User{
		Username: "alice_001",
		Email:    &email,
		Role:     auth.UserRoleMember,
		Status:   auth.UserStatusActive,
	}

	if err := store.Users.Create(ctx, user); err != nil {
		t.Fatalf("expected create user to succeed, got %v", err)
	}

	gotByUsername, err := store.Users.GetByUsername(ctx, "alice_001")
	if err != nil {
		t.Fatalf("expected get by username to succeed, got %v", err)
	}
	if gotByUsername.ID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, gotByUsername.ID)
	}

	gotByEmail, err := store.Users.GetByEmail(ctx, "ALICE@example.com")
	if err != nil {
		t.Fatalf("expected get by email to succeed, got %v", err)
	}
	if gotByEmail.ID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, gotByEmail.ID)
	}

	now := time.Now().UTC()
	if err := store.Users.TouchLogin(ctx, user.ID, now); err != nil {
		t.Fatalf("expected touch login to succeed, got %v", err)
	}

	dupUser := &entity.User{
		Username: "alice_001",
		Role:     auth.UserRoleMember,
		Status:   auth.UserStatusActive,
	}
	if err := store.Users.Create(ctx, dupUser); !repo.IsUniqueViolation(err) {
		t.Fatalf("expected unique violation, got %v", err)
	}
}
