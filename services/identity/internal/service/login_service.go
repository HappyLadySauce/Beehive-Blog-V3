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

// LoginService handles local login use cases.
// LoginService 处理本地登录用例。
type LoginService struct {
	deps Dependencies
}

// NewLoginService creates a LoginService instance.
// NewLoginService 创建 LoginService 实例。
func NewLoginService(deps Dependencies) *LoginService {
	return &LoginService{deps: deps}
}

// Execute authenticates a local user and creates a new session.
// Execute 认证本地用户并创建新会话。
func (s *LoginService) Execute(ctx context.Context, in LoginLocalUserInput) (*AuthResult, error) {
	identifier, err := auth.NormalizeLoginIdentifier(in.LoginIdentifier)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, err.Error())
	}
	if strings.TrimSpace(in.Password) == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "password is required")
	}

	user, err := s.deps.Store.Users.FindByLoginIdentifier(ctx, identifier)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeIdentityInvalidCredentials, "invalid credentials")
		}
		return nil, err
	}
	if err := validateActiveUserStatus(user.Status); err != nil {
		return nil, err
	}

	credential, err := s.deps.Store.CredentialLocals.GetByUserID(ctx, user.ID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeIdentityInvalidCredentials, "invalid credentials")
		}
		return nil, err
	}
	if err := auth.VerifyPassword(credential.PasswordHash, in.Password); err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "invalid credentials")
	}

	now := s.deps.Clock()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	authSource := auth.AuthSourceLocal

	var createdSession *entity.UserSession
	err = withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		if err := store.Users.TouchLogin(ctx, user.ID, now); err != nil {
			return err
		}

		createdSession = &entity.UserSession{
			UserID:     user.ID,
			AuthSource: auth.AuthSourceLocal,
			ClientType: stringPtr(in.ClientType),
			DeviceID:   stringPtr(in.DeviceID),
			DeviceName: stringPtr(in.DeviceName),
			IPAddress:  stringPtr(in.ClientIP),
			UserAgent:  stringPtr(in.UserAgent),
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(s.deps.Config.Security.RefreshTokenTTLSeconds) * time.Second),
		}
		if err := store.UserSessions.Create(ctx, createdSession); err != nil {
			return err
		}

		refreshRow := &entity.RefreshToken{
			SessionID: createdSession.ID,
			TokenHash: refreshTokenHash,
			IssuedAt:  now,
			ExpiresAt: createdSession.ExpiresAt,
		}
		if err := store.RefreshTokens.Create(ctx, refreshRow); err != nil {
			return err
		}

		writeAudit(ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &createdSession.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventLoginLocal,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(in.ClientIP),
			UserAgent:  stringPtr(in.UserAgent),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"login_identifier": identifier,
			}),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	user.LastLoginAt = &now
	accessToken, accessExpiresAt, err := issueAccessToken(
		s.deps.Config.Security.AccessTokenSecret,
		s.deps.Config.Security.AccessTokenTTLSeconds,
		user,
		createdSession,
		now,
	)
	if err != nil {
		return nil, err
	}

	return buildAuthResult(user, createdSession, accessToken, refreshToken, accessExpiresAt), nil
}
