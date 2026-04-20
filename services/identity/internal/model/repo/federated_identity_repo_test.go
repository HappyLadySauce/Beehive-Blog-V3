package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestFederatedIdentityRepository verifies federated identity persistence.
// TestFederatedIdentityRepository 验证联邦身份持久化。
func TestFederatedIdentityRepository(t *testing.T) {
	t.Parallel()

	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	ctx := context.Background()

	fed := testkit.CreateFederatedIdentity(t, store, user.ID, "github", "12345")

	got, err := store.FederatedIdentities.GetByProviderSubject(ctx, "github", "12345")
	if err != nil {
		t.Fatalf("expected federated identity lookup to succeed, got %v", err)
	}
	if got.ID != fed.ID {
		t.Fatalf("expected federated identity id %d, got %d", fed.ID, got.ID)
	}

	lastLoginAt := time.Now().UTC()
	displayName := "Octocat"
	if err := store.FederatedIdentities.TouchLogin(ctx, fed.ID, nil, &displayName, nil, []byte(`{"id":12345}`), lastLoginAt); err != nil {
		t.Fatalf("expected touch login to succeed, got %v", err)
	}
}
