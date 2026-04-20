package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestLogoutServiceExecute verifies session revocation behavior.
// TestLogoutServiceExecute 验证会话吊销行为。
func TestLogoutServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("logout success", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID)
		token := testkit.CreateRefreshToken(t, deps.Store, session.ID, auth.HashRefreshToken("refresh-token"))

		svc := service.NewLogoutService(deps)
		if err := svc.Execute(context.Background(), service.LogoutSessionInput{
			SessionID: session.ID,
			ClientIP:  "127.0.0.1",
		}); err != nil {
			t.Fatalf("expected logout to succeed, got %v", err)
		}

		updatedSession, err := deps.Store.UserSessions.GetByID(context.Background(), session.ID)
		if err != nil {
			t.Fatalf("expected session lookup to succeed, got %v", err)
		}
		if updatedSession.Status != auth.SessionStatusRevoked {
			t.Fatalf("expected revoked session, got %s", updatedSession.Status)
		}

		refreshRow, err := deps.Store.RefreshTokens.GetActiveForUpdateByHash(context.Background(), token.TokenHash)
		if err == nil || refreshRow != nil {
			t.Fatalf("expected refresh token to be revoked")
		}
	})

	t.Run("session not found", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewLogoutService(deps)

		err := svc.Execute(context.Background(), service.LogoutSessionInput{SessionID: 9999})
		if !errors.Is(err, errs.E(errs.CodeIdentitySessionNotFound)) {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}
