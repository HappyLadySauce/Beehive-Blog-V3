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

// RegisterService handles local registration use cases.
// RegisterService 处理本地注册用例。
type RegisterService struct {
	deps Dependencies
}

// NewRegisterService creates a RegisterService instance.
// NewRegisterService 创建 RegisterService 实例。
func NewRegisterService(deps Dependencies) *RegisterService {
	return &RegisterService{deps: deps}
}

// Execute registers a local account and creates the initial session.
// Execute 注册本地账号并创建初始会话。
func (s *RegisterService) Execute(ctx context.Context, in RegisterLocalUserInput) (*AuthResult, error) {
	username, err := auth.NormalizeUsername(in.Username)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, err.Error())
	}
	email, err := auth.NormalizeEmail(in.Email)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, err.Error())
	}
	nickname, err := auth.NormalizeNickname(in.Nickname)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, err.Error())
	}
	if err := auth.ValidatePassword(in.Password); err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, err.Error())
	}

	if _, err := s.deps.Store.Users.GetByUsername(ctx, username); err == nil {
		return nil, errs.New(errs.CodeIdentityUsernameAlreadyExists, "username already exists")
	} else if !repo.IsNotFound(err) {
		return nil, err
	}
	if email != "" {
		if _, err := s.deps.Store.Users.GetByEmail(ctx, email); err == nil {
			return nil, errs.New(errs.CodeIdentityEmailAlreadyExists, "email already exists")
		} else if !repo.IsNotFound(err) {
			return nil, err
		}
	}

	hashedPassword, err := auth.HashPassword(in.Password, s.deps.Config.Security.PasswordHashCost)
	if err != nil {
		return nil, err
	}

	now := s.deps.Clock()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	authSource := auth.AuthSourceLocal

	var createdUser *entity.User
	var createdSession *entity.UserSession

	err = withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		createdUser = &entity.User{
			Username:    username,
			Email:       stringPtr(email),
			Nickname:    stringPtr(strings.TrimSpace(nickname)),
			Role:        auth.UserRoleMember,
			Status:      auth.UserStatusActive,
			LastLoginAt: &now,
		}
		if err := store.Users.Create(ctx, createdUser); err != nil {
			return err
		}

		credential := &entity.CredentialLocal{
			UserID:            createdUser.ID,
			PasswordHash:      hashedPassword,
			PasswordUpdatedAt: now,
		}
		if err := store.CredentialLocals.Create(ctx, credential); err != nil {
			return err
		}

		createdSession = &entity.UserSession{
			UserID:     createdUser.ID,
			AuthSource: auth.AuthSourceLocal,
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(s.deps.Config.Security.RefreshTokenTTLSeconds) * time.Second),
			IPAddress:  stringPtr(in.ClientIP),
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
			UserID:     &createdUser.ID,
			SessionID:  &createdSession.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventRegisterLocal,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(in.ClientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"username": createdUser.Username,
			}),
		})

		return nil
	})
	if err != nil {
		if repo.IsUniqueViolation(err) {
			return nil, errs.Wrap(err, errs.CodeIdentityEmailAlreadyExists, "username or email already exists")
		}
		return nil, err
	}

	accessToken, accessExpiresAt, err := issueAccessToken(
		s.deps.Config.Security.AccessTokenSecret,
		s.deps.Config.Security.AccessTokenTTLSeconds,
		createdUser,
		createdSession,
		now,
	)
	if err != nil {
		return nil, err
	}

	return buildAuthResult(createdUser, createdSession, accessToken, refreshToken, accessExpiresAt), nil
}
