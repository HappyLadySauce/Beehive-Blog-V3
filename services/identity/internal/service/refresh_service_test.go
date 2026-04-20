package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestRefreshServiceExecute verifies refresh token rotation behavior.
// TestRefreshServiceExecute 验证 refresh token 轮换行为。
func TestRefreshServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("refresh success", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID)
		rawToken := "refresh-token"
		testkit.CreateRefreshToken(t, deps.Store, session.ID, auth.HashRefreshToken(rawToken))

		svc := service.NewRefreshService(deps)
		result, err := svc.Execute(context.Background(), service.RefreshSessionTokenInput{
			RefreshToken: rawToken,
			ClientIP:     "127.0.0.1",
		})
		if err != nil {
			t.Fatalf("expected refresh to succeed, got %v", err)
		}
		if result.RefreshToken == rawToken {
			t.Fatalf("expected refresh token rotation to produce a new token")
		}
	})

	t.Run("revoked session fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID, func(s *entity.UserSession) {
			s.Status = auth.SessionStatusRevoked
		})
		rawToken := "refresh-token"
		testkit.CreateRefreshToken(t, deps.Store, session.ID, auth.HashRefreshToken(rawToken))

		svc := service.NewRefreshService(deps)
		_, err := svc.Execute(context.Background(), service.RefreshSessionTokenInput{
			RefreshToken: rawToken,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySessionRevoked)) {
			t.Fatalf("expected session revoked error, got %v", err)
		}
	})

	t.Run("expired refresh token fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID)
		rawToken := "refresh-token"
		testkit.CreateRefreshToken(t, deps.Store, session.ID, auth.HashRefreshToken(rawToken), func(token *entity.RefreshToken) {
			token.ExpiresAt = now.Add(-time.Minute)
		})

		svc := service.NewRefreshService(deps)
		_, err := svc.Execute(context.Background(), service.RefreshSessionTokenInput{
			RefreshToken: rawToken,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentityRefreshTokenExpired)) {
			t.Fatalf("expected refresh token expired error, got %v", err)
		}
	})

	t.Run("disabled user fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Status = auth.UserStatusDisabled
		})
		session := testkit.CreateSession(t, deps.Store, user.ID)
		rawToken := "refresh-token"
		testkit.CreateRefreshToken(t, deps.Store, session.ID, auth.HashRefreshToken(rawToken))

		svc := service.NewRefreshService(deps)
		_, err := svc.Execute(context.Background(), service.RefreshSessionTokenInput{
			RefreshToken: rawToken,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentityAccountDisabled)) {
			t.Fatalf("expected account disabled error, got %v", err)
		}
	})
}
