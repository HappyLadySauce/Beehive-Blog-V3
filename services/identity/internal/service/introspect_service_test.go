package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestIntrospectServiceExecute verifies active and inactive token states.
// TestIntrospectServiceExecute 验证 token 的 active 与 inactive 状态。
func TestIntrospectServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("active token", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID)
		token, _, err := auth.IssueAccessToken(
			deps.Config.Security.AccessTokenSecret,
			15*time.Minute,
			user.ID,
			user.Role,
			user.Status,
			session.ID,
			session.AuthSource,
			now,
		)
		if err != nil {
			t.Fatalf("expected issue token to succeed, got %v", err)
		}

		svc := service.NewIntrospectService(deps)
		result, err := svc.Execute(context.Background(), service.IntrospectAccessTokenInput{AccessToken: token})
		if err != nil {
			t.Fatalf("expected introspect to succeed, got %v", err)
		}
		if !result.Active {
			t.Fatalf("expected token to be active")
		}
	})

	t.Run("disabled user inactive", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Status = auth.UserStatusDisabled
		})
		session := testkit.CreateSession(t, deps.Store, user.ID)
		token, _, err := auth.IssueAccessToken(
			deps.Config.Security.AccessTokenSecret,
			15*time.Minute,
			user.ID,
			user.Role,
			user.Status,
			session.ID,
			session.AuthSource,
			now,
		)
		if err != nil {
			t.Fatalf("expected issue token to succeed, got %v", err)
		}

		svc := service.NewIntrospectService(deps)
		result, err := svc.Execute(context.Background(), service.IntrospectAccessTokenInput{AccessToken: token})
		if err != nil {
			t.Fatalf("expected introspect to succeed, got %v", err)
		}
		if result.Active {
			t.Fatalf("expected disabled user token to be inactive")
		}
	})

	t.Run("expired token inactive", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID)
		token, _, err := auth.IssueAccessToken(
			deps.Config.Security.AccessTokenSecret,
			-time.Minute,
			user.ID,
			user.Role,
			user.Status,
			session.ID,
			session.AuthSource,
			now,
		)
		if err != nil {
			t.Fatalf("expected issue token to succeed, got %v", err)
		}

		svc := service.NewIntrospectService(deps)
		result, err := svc.Execute(context.Background(), service.IntrospectAccessTokenInput{AccessToken: token})
		if err != nil {
			t.Fatalf("expected introspect to succeed, got %v", err)
		}
		if result.Active {
			t.Fatalf("expected expired token to be inactive")
		}
	})

	t.Run("revoked session inactive", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		session := testkit.CreateSession(t, deps.Store, user.ID, func(s *entity.UserSession) {
			s.Status = auth.SessionStatusRevoked
		})
		token, _, err := auth.IssueAccessToken(
			deps.Config.Security.AccessTokenSecret,
			15*time.Minute,
			user.ID,
			user.Role,
			user.Status,
			session.ID,
			session.AuthSource,
			now,
		)
		if err != nil {
			t.Fatalf("expected issue token to succeed, got %v", err)
		}

		svc := service.NewIntrospectService(deps)
		result, err := svc.Execute(context.Background(), service.IntrospectAccessTokenInput{AccessToken: token})
		if err != nil {
			t.Fatalf("expected introspect to succeed, got %v", err)
		}
		if result.Active {
			t.Fatalf("expected revoked session token to be inactive")
		}
	})
}
