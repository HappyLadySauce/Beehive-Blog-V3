package repo_test

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestIdentityAuditRepository verifies audit persistence.
// TestIdentityAuditRepository 验证审计记录持久化。
func TestIdentityAuditRepository(t *testing.T) {
	t.Parallel()

	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	session := testkit.CreateSession(t, store, user.ID)

	record := &entity.IdentityAudit{
		UserID:     &user.ID,
		SessionID:  &session.ID,
		AuthSource: testStringPtr(auth.AuthSourceLocal),
		EventType:  auth.AuditEventLoginLocal,
		Result:     auth.AuditResultSuccess,
	}
	if err := store.IdentityAudits.Create(context.Background(), record); err != nil {
		t.Fatalf("expected audit create to succeed, got %v", err)
	}
	if record.ID == 0 {
		t.Fatalf("expected audit id to be populated")
	}
}

func testStringPtr(value string) *string {
	return &value
}
