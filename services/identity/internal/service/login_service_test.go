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

// TestLoginServiceExecute verifies local login behavior.
// TestLoginServiceExecute 验证本地登录行为。
func TestLoginServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("login by username", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Username = "alice_001"
		})
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		result, err := svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: "alice_001",
			Password:        "password123",
		})
		if err != nil {
			t.Fatalf("expected login to succeed, got %v", err)
		}
		if result.User.ID != user.ID {
			t.Fatalf("expected user_id=%d, got %d", user.ID, result.User.ID)
		}
	})

	t.Run("login by email", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		result, err := svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: *user.Email,
			Password:        "password123",
		})
		if err != nil {
			t.Fatalf("expected login by email to succeed, got %v", err)
		}
		if result.User.ID != user.ID {
			t.Fatalf("expected user_id=%d, got %d", user.ID, result.User.ID)
		}
	})

	t.Run("disabled user fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Status = auth.UserStatusDisabled
		})
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		_, err = svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: user.Username,
			Password:        "password123",
		})
		if !service.IsKind(err, service.ErrorKindFailedPrecondition) {
			t.Fatalf("expected failed precondition error, got %v", err)
		}
	})

	t.Run("password mismatch fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		_, err = svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: user.Username,
			Password:        "wrong-password",
		})
		if !service.IsKind(err, service.ErrorKindUnauthenticated) {
			t.Fatalf("expected unauthenticated error, got %v", err)
		}
	})

	t.Run("pending user fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Status = auth.UserStatusPending
		})
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		_, err = svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: user.Username,
			Password:        "password123",
		})
		if !service.IsKind(err, service.ErrorKindFailedPrecondition) {
			t.Fatalf("expected failed precondition error, got %v", err)
		}
	})

	t.Run("locked user fails", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store, func(u *entity.User) {
			u.Status = auth.UserStatusLocked
		})
		hash, err := auth.HashPassword("password123", deps.Config.Security.PasswordHashCost)
		if err != nil {
			t.Fatalf("expected hash password to succeed, got %v", err)
		}
		testkit.CreateCredentialLocal(t, deps.Store, user.ID, hash)

		svc := service.NewLoginService(deps)
		_, err = svc.Execute(context.Background(), service.LoginLocalUserInput{
			LoginIdentifier: user.Username,
			Password:        "password123",
		})
		if !service.IsKind(err, service.ErrorKindFailedPrecondition) {
			t.Fatalf("expected failed precondition error, got %v", err)
		}
	})
}
