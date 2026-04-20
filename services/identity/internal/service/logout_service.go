package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// LogoutService handles session revocation.
// LogoutService 处理会话吊销。
type LogoutService struct {
	deps Dependencies
}

// NewLogoutService creates a LogoutService instance.
// NewLogoutService 创建 LogoutService 实例。
func NewLogoutService(deps Dependencies) *LogoutService {
	return &LogoutService{deps: deps}
}

// Execute revokes the target session and its refresh tokens.
// Execute 吊销目标会话及其 refresh tokens。
func (s *LogoutService) Execute(ctx context.Context, in LogoutSessionInput) error {
	if in.SessionID <= 0 {
		return errs.New(errs.CodeIdentityInvalidArgument, "session_id is invalid")
	}

	now := s.deps.Clock()
	return withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		session, err := store.UserSessions.GetForUpdateByID(ctx, in.SessionID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentitySessionNotFound, "session not found")
			}
			return err
		}

		if session.Status == auth.SessionStatusActive {
			if err := store.UserSessions.Revoke(ctx, session.ID, now); err != nil {
				return err
			}
		}
		if err := store.RefreshTokens.RevokeActiveBySessionID(ctx, session.ID, now); err != nil {
			return err
		}

		authSource := session.AuthSource
		writeAudit(ctx, store, auditInput{
			UserID:     &session.UserID,
			SessionID:  &session.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventLogoutSession,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(in.ClientIP),
		})

		return nil
	})
}
