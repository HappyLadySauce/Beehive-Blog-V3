package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// RefreshService handles refresh token rotation.
// RefreshService 处理 refresh token 轮换。
type RefreshService struct {
	deps Dependencies
}

// NewRefreshService creates a RefreshService instance.
// NewRefreshService 创建 RefreshService 实例。
func NewRefreshService(deps Dependencies) *RefreshService {
	return &RefreshService{deps: deps}
}

// Execute rotates the refresh token and reissues the access token.
// Execute 轮换 refresh token 并重新签发 access token。
func (s *RefreshService) Execute(ctx context.Context, in RefreshSessionTokenInput) (*AuthResult, error) {
	rawRefreshToken := strings.TrimSpace(in.RefreshToken)
	if rawRefreshToken == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "refresh_token is required")
	}

	now := s.deps.Clock()
	refreshTokenHash := auth.HashRefreshToken(rawRefreshToken)
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	newRefreshTokenHash := auth.HashRefreshToken(newRefreshToken)

	var user *entity.User
	var session *entity.UserSession

	err = withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		currentToken, err := store.RefreshTokens.GetActiveForUpdateByHash(ctx, refreshTokenHash)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityInvalidRefreshToken, "invalid refresh token")
			}
			return err
		}
		if currentToken.ExpiresAt.Before(now) {
			_ = store.RefreshTokens.Revoke(ctx, currentToken.ID, now)
			return errs.New(errs.CodeIdentityRefreshTokenExpired, "refresh token expired")
		}

		session, err = store.UserSessions.GetForUpdateByID(ctx, currentToken.SessionID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentitySessionRevoked, "session revoked")
			}
			return err
		}
		if session.Status != auth.SessionStatusActive {
			return errs.New(errs.CodeIdentitySessionRevoked, "session revoked")
		}
		if session.ExpiresAt.Before(now) {
			_ = store.UserSessions.MarkExpired(ctx, session.ID, now)
			return errs.New(errs.CodeIdentityRefreshTokenExpired, "refresh token expired")
		}

		user, err = store.Users.GetByID(ctx, session.UserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityAccountNotFound, "user not found")
			}
			return err
		}
		if err := validateActiveUserStatus(user.Status); err != nil {
			return err
		}

		session.LastSeenAt = &now
		session.ExpiresAt = now.Add(time.Duration(s.deps.Config.Security.RefreshTokenTTLSeconds) * time.Second)
		if err := store.UserSessions.TouchActive(ctx, session.ID, now, session.ExpiresAt); err != nil {
			return err
		}
		if err := store.RefreshTokens.Revoke(ctx, currentToken.ID, now); err != nil {
			return err
		}

		newTokenRow := &entity.RefreshToken{
			SessionID:          session.ID,
			TokenHash:          newRefreshTokenHash,
			IssuedAt:           now,
			ExpiresAt:          session.ExpiresAt,
			RotatedFromTokenID: &currentToken.ID,
		}
		if err := store.RefreshTokens.Create(ctx, newTokenRow); err != nil {
			return err
		}

		authSource := session.AuthSource
		writeAudit(ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &session.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventRefreshSession,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(in.ClientIP),
			UserAgent:  stringPtr(in.UserAgent),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, accessExpiresAt, err := issueAccessToken(
		s.deps.Config.Security.AccessTokenSecret,
		s.deps.Config.Security.AccessTokenTTLSeconds,
		user,
		session,
		now,
	)
	if err != nil {
		return nil, err
	}

	return buildAuthResult(user, session, accessToken, newRefreshToken, accessExpiresAt), nil
}
