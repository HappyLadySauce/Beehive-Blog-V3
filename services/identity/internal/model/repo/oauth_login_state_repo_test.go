package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestOAuthLoginStateRepository verifies oauth state persistence and consume behavior.
// TestOAuthLoginStateRepository 验证 oauth state 持久化与消费行为。
func TestOAuthLoginStateRepository(t *testing.T) {
	t.Parallel()

	store := testkit.NewStore(t)
	ctx := context.Background()

	row := &entity.OAuthLoginState{
		Provider:    "github",
		State:       "state-1",
		RedirectURI: "https://example.com/callback",
		ExpiresAt:   time.Now().UTC().Add(time.Minute),
	}
	if err := store.OAuthLoginStates.Create(ctx, row); err != nil {
		t.Fatalf("expected create oauth state to succeed, got %v", err)
	}

	got, err := store.OAuthLoginStates.GetForUpdateByProviderState(ctx, "github", "state-1")
	if err != nil {
		t.Fatalf("expected get oauth state to succeed, got %v", err)
	}
	if got.ID != row.ID {
		t.Fatalf("expected state id %d, got %d", row.ID, got.ID)
	}

	consumedAt := time.Now().UTC()
	if err := store.OAuthLoginStates.Consume(ctx, row.ID, consumedAt); err != nil {
		t.Fatalf("expected consume oauth state to succeed, got %v", err)
	}
}
