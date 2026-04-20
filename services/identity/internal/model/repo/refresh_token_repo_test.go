package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestRefreshTokenRepository verifies refresh token persistence and revocation.
// TestRefreshTokenRepository 验证 refresh token 持久化与吊销行为。
func TestRefreshTokenRepository(t *testing.T) {
	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	session := testkit.CreateSession(t, store, user.ID)
	ctx := context.Background()

	row := testkit.CreateRefreshToken(t, store, session.ID, auth.HashRefreshToken("refresh-token"))

	got, err := store.RefreshTokens.GetActiveForUpdateByHash(ctx, row.TokenHash)
	if err != nil {
		t.Fatalf("expected get refresh token to succeed, got %v", err)
	}
	if got.ID != row.ID {
		t.Fatalf("expected refresh token id %d, got %d", row.ID, got.ID)
	}

	revokedAt := time.Now().UTC()
	if err := store.RefreshTokens.Revoke(ctx, row.ID, revokedAt); err != nil {
		t.Fatalf("expected revoke refresh token to succeed, got %v", err)
	}
	if _, err := store.RefreshTokens.GetActiveForUpdateByHash(ctx, row.TokenHash); err == nil {
		t.Fatalf("expected revoked refresh token lookup to fail")
	}

	activeRow := testkit.CreateRefreshToken(t, store, session.ID, auth.HashRefreshToken("refresh-token-2"))
	if err := store.RefreshTokens.RevokeActiveBySessionID(ctx, session.ID, revokedAt); err != nil {
		t.Fatalf("expected revoke active by session to succeed, got %v", err)
	}
	if _, err := store.RefreshTokens.GetActiveForUpdateByHash(ctx, activeRow.TokenHash); err == nil {
		t.Fatalf("expected session refresh tokens to be revoked")
	}
}
