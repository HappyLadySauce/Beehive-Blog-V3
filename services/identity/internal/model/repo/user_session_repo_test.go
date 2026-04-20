package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestUserSessionRepository verifies session persistence and state transitions.
// TestUserSessionRepository 验证会话持久化与状态流转。
func TestUserSessionRepository(t *testing.T) {
	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	ctx := context.Background()
	now := time.Now().UTC()

	session := &entity.UserSession{
		UserID:     user.ID,
		AuthSource: auth.AuthSourceLocal,
		Status:     auth.SessionStatusActive,
		LastSeenAt: &now,
		ExpiresAt:  now.Add(time.Hour),
	}
	if err := store.UserSessions.Create(ctx, session); err != nil {
		t.Fatalf("expected create session to succeed, got %v", err)
	}

	got, err := store.UserSessions.GetForUpdateByID(ctx, session.ID)
	if err != nil {
		t.Fatalf("expected get session to succeed, got %v", err)
	}
	if got.ID != session.ID {
		t.Fatalf("expected session id %d, got %d", session.ID, got.ID)
	}

	touchAt := now.Add(10 * time.Minute)
	expiresAt := now.Add(2 * time.Hour)
	if err := store.UserSessions.TouchActive(ctx, session.ID, touchAt, expiresAt); err != nil {
		t.Fatalf("expected touch active to succeed, got %v", err)
	}

	if err := store.UserSessions.Revoke(ctx, session.ID, touchAt); err != nil {
		t.Fatalf("expected revoke to succeed, got %v", err)
	}

	revoked, err := store.UserSessions.GetByID(ctx, session.ID)
	if err != nil {
		t.Fatalf("expected get revoked session to succeed, got %v", err)
	}
	if revoked.Status != auth.SessionStatusRevoked {
		t.Fatalf("expected revoked status, got %s", revoked.Status)
	}

	expiredSession := &entity.UserSession{
		UserID:     user.ID,
		AuthSource: auth.AuthSourceLocal,
		Status:     auth.SessionStatusActive,
		LastSeenAt: &now,
		ExpiresAt:  now.Add(time.Hour),
	}
	if err := store.UserSessions.Create(ctx, expiredSession); err != nil {
		t.Fatalf("expected create second session to succeed, got %v", err)
	}
	if err := store.UserSessions.MarkExpired(ctx, expiredSession.ID, touchAt); err != nil {
		t.Fatalf("expected mark expired to succeed, got %v", err)
	}

	expired, err := store.UserSessions.GetByID(ctx, expiredSession.ID)
	if err != nil {
		t.Fatalf("expected get expired session to succeed, got %v", err)
	}
	if expired.Status != auth.SessionStatusExpired {
		t.Fatalf("expected expired status, got %s", expired.Status)
	}
}
