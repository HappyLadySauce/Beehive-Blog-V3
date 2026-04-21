package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityprovider "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

type validatedOAuthState struct {
	ID int64
}

type ssoStateFailure struct {
	reason string
	err    error
}

func (e *ssoStateFailure) Error() string {
	if e == nil || e.err == nil {
		return ""
	}

	return e.err.Error()
}

func (e *ssoStateFailure) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.err
}

// SSOFinishService handles SSO callback completion.
// SSOFinishService 处理 SSO callback 完成流程。
type SSOFinishService struct {
	deps Dependencies
}

// NewSSOFinishService creates an SSOFinishService instance.
// NewSSOFinishService 创建 SSOFinishService 实例。
func NewSSOFinishService(deps Dependencies) *SSOFinishService {
	return &SSOFinishService{deps: deps}
}

// Execute completes the SSO callback flow for fully implemented providers.
// Execute 为完整实现的 provider 完成 SSO callback 流程。
func (s *SSOFinishService) Execute(ctx context.Context, in FinishSSOInput) (*AuthResult, error) {
	providerName, err := auth.NormalizeProvider(in.Provider)
	if err != nil {
		s.writeFailureAudit(ctx, in, strings.TrimSpace(in.Provider), "invalid_provider")
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, "unsupported provider")
	}

	callbackProvider, ok := s.deps.Providers.GetCallback(providerName)
	if !ok {
		s.writeFailureAudit(ctx, in, providerName, "provider_not_ready")
		return nil, errs.New(errs.CodeIdentitySSOProviderNotReady, "sso provider is not ready")
	}
	if !callbackProvider.Enabled() {
		s.writeFailureAudit(ctx, in, providerName, "provider_disabled")
		return nil, errs.New(errs.CodeIdentitySSOProviderDisabled, "sso provider is disabled")
	}
	if !callbackProvider.LoginReady() {
		s.writeFailureAudit(ctx, in, providerName, "provider_not_ready")
		return nil, errs.New(errs.CodeIdentitySSOProviderNotReady, "sso provider is not ready")
	}
	if strings.TrimSpace(in.Code) == "" || strings.TrimSpace(in.State) == "" {
		s.writeFailureAudit(ctx, in, providerName, "missing_code_or_state")
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "code and state are required")
	}

	redirectURI := strings.TrimSpace(in.RedirectURI)
	if redirectURI == "" || redirectURI != strings.TrimSpace(callbackProvider.RedirectURL()) {
		s.writeFailureAudit(ctx, in, providerName, "redirect_uri_mismatch")
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "redirect_uri does not match configured provider redirect")
	}

	stateRow, err := s.validateOAuthState(ctx, providerName, strings.TrimSpace(in.State), redirectURI)
	if err != nil {
		s.writeFailureAudit(ctx, in, providerName, ssoStateFailureReason(err, "invalid_state"))
		return nil, err
	}

	providerAccessToken, err := callbackProvider.ExchangeCode(ctx, in.Code, redirectURI)
	if err != nil {
		s.writeFailureAudit(ctx, in, providerName, "exchange_code_failed")
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "exchange provider code failed")
	}
	profile, _, err := callbackProvider.FetchProfile(ctx, providerAccessToken)
	if err != nil {
		s.writeFailureAudit(ctx, in, providerName, "profile_fetch_failed")
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "fetch provider profile failed")
	}

	now := s.deps.Clock()
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)

	var user *entity.User
	var session *entity.UserSession

	err = withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		lockedState, err := store.OAuthLoginStates.GetForUpdateByProviderState(ctx, providerName, strings.TrimSpace(in.State))
		if err != nil {
			if repo.IsNotFound(err) {
				return newSSOStateFailure("state_not_found", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
			}
			return err
		}
		if lockedState.ID != stateRow.ID {
			return newSSOStateFailure("state_already_consumed", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if lockedState.ConsumedAt != nil {
			return newSSOStateFailure("state_already_consumed", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if lockedState.ExpiresAt.Before(now) {
			return newSSOStateFailure("state_expired", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if lockedState.RedirectURI != redirectURI {
			return newSSOStateFailure("state_redirect_mismatch", errs.CodeIdentityInvalidArgument, "redirect_uri mismatch")
		}

		fed, err := store.FederatedIdentities.GetByProviderSubject(ctx, providerName, profile.Subject)
		if err != nil && !repo.IsNotFound(err) {
			return err
		}

		if fed != nil {
			user, err = store.Users.GetByID(ctx, fed.UserID)
			if err != nil {
				return err
			}
			if err := validateActiveUserStatus(user.Status); err != nil {
				return err
			}
			if err := store.FederatedIdentities.TouchLogin(
				ctx,
				fed.ID,
				profile.Email,
				stringPtr(profile.DisplayName),
				profile.AvatarURL,
				profile.RawProfile,
				now,
			); err != nil {
				return err
			}
		} else {
			if profile.Email != nil && *profile.Email != "" {
				user, err = store.Users.GetByEmail(ctx, *profile.Email)
				if err != nil && !repo.IsNotFound(err) {
					return err
				}
			}
			if user == nil {
				username, err := buildUniqueGitHubUsername(ctx, store, profile.Login)
				if err != nil {
					return err
				}
				nickname := profile.DisplayName
				if strings.TrimSpace(nickname) == "" {
					nickname = profile.Login
				}
				user = &entity.User{
					Username:    username,
					Email:       profile.Email,
					Nickname:    stringPtr(nickname),
					AvatarURL:   profile.AvatarURL,
					Role:        auth.UserRoleMember,
					Status:      auth.UserStatusActive,
					LastLoginAt: &now,
				}
				if err := store.Users.Create(ctx, user); err != nil {
					return err
				}
			} else {
				if err := validateActiveUserStatus(user.Status); err != nil {
					return err
				}
			}

			fed = &entity.FederatedIdentity{
				UserID:              user.ID,
				Provider:            providerName,
				ProviderSubject:     profile.Subject,
				ProviderSubjectType: profile.SubjectType,
				ProviderLogin:       stringPtr(profile.Login),
				ProviderEmail:       profile.Email,
				ProviderDisplayName: stringPtr(profile.DisplayName),
				AvatarURL:           profile.AvatarURL,
				AppIDOrClientID:     profile.ProviderClientID,
				AccessScope:         profile.RequestedScopes,
				RawProfile:          profile.RawProfile,
				LastLoginAt:         &now,
			}
			if err := store.FederatedIdentities.Create(ctx, fed); err != nil {
				return err
			}
		}

		if err := store.Users.TouchLogin(ctx, user.ID, now); err != nil {
			return err
		}

		session = &entity.UserSession{
			UserID:     user.ID,
			AuthSource: auth.AuthSourceSSO,
			ClientType: stringPtr(in.ClientType),
			DeviceID:   stringPtr(in.DeviceID),
			DeviceName: stringPtr(in.DeviceName),
			IPAddress:  stringPtr(in.ClientIP),
			UserAgent:  stringPtr(in.UserAgent),
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(s.deps.Config.Security.RefreshTokenTTLSeconds) * time.Second),
		}
		if err := store.UserSessions.Create(ctx, session); err != nil {
			return err
		}

		refreshRow := &entity.RefreshToken{
			SessionID: session.ID,
			TokenHash: refreshTokenHash,
			IssuedAt:  now,
			ExpiresAt: session.ExpiresAt,
		}
		if err := store.RefreshTokens.Create(ctx, refreshRow); err != nil {
			return err
		}
		if err := store.OAuthLoginStates.Consume(ctx, lockedState.ID, now); err != nil {
			return err
		}

		authSource := auth.AuthSourceSSO
		writeAudit(ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &session.ID,
			Provider:   stringPtr(providerName),
			AuthSource: &authSource,
			EventType:  auth.AuditEventFinishSSO,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(in.ClientIP),
			UserAgent:  stringPtr(in.UserAgent),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"provider_subject": profile.Subject,
			}),
		})

		return nil
	})
	if err != nil {
		if shouldWriteSSOStateFailureAudit(err) {
			s.writeFailureAudit(ctx, in, providerName, ssoStateFailureReason(err, "invalid_state"))
		}
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

	return buildAuthResult(user, session, accessToken, refreshToken, accessExpiresAt), nil
}

// validateOAuthState validates and locks the oauth state in a short transaction before upstream calls.
// validateOAuthState 在调用上游前，通过短事务校验并锁定 oauth state。
func (s *SSOFinishService) validateOAuthState(ctx context.Context, providerName, state, redirectURI string) (*validatedOAuthState, error) {
	now := s.deps.Clock()
	result := &validatedOAuthState{}

	err := withTransaction(ctx, s.deps.Store, func(store *repo.Store) error {
		stateRow, err := store.OAuthLoginStates.GetForUpdateByProviderState(ctx, providerName, state)
		if err != nil {
			if repo.IsNotFound(err) {
				return newSSOStateFailure("state_not_found", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
			}
			return err
		}
		if stateRow.ConsumedAt != nil {
			return newSSOStateFailure("state_already_consumed", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if stateRow.ExpiresAt.Before(now) {
			return newSSOStateFailure("state_expired", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if stateRow.RedirectURI != redirectURI {
			return newSSOStateFailure("state_redirect_mismatch", errs.CodeIdentityInvalidArgument, "redirect_uri mismatch")
		}

		result.ID = stateRow.ID
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// writeFailureAudit records SSO callback failures on a best-effort basis.
// writeFailureAudit 以尽力而为方式记录 SSO callback 失败审计。
func (s *SSOFinishService) writeFailureAudit(ctx context.Context, in FinishSSOInput, providerName, reason string) {
	writeAudit(ctx, s.deps.Store, auditInput{
		Provider:  stringPtr(providerName),
		EventType: auth.AuditEventFinishSSO,
		Result:    auth.AuditResultFailure,
		ClientIP:  stringPtr(in.ClientIP),
		UserAgent: stringPtr(in.UserAgent),
		Detail: auth.MarshalAuditDetail(map[string]any{
			"reason": reason,
		}),
	})
}

func newSSOStateFailure(reason string, code errs.Code, message string) error {
	return &ssoStateFailure{
		reason: reason,
		err:    errs.New(code, message),
	}
}

func ssoStateFailureReason(err error, fallback string) string {
	var stateErr *ssoStateFailure
	if errors.As(err, &stateErr) && strings.TrimSpace(stateErr.reason) != "" {
		return stateErr.reason
	}

	return fallback
}

func shouldWriteSSOStateFailureAudit(err error) bool {
	var stateErr *ssoStateFailure
	return errors.As(err, &stateErr)
}

// errorsIsSSOStateInvalid reports whether the error represents an invalid or consumed SSO state.
// errorsIsSSOStateInvalid 判断错误是否表示 SSO state 无效或已消费。
func errorsIsSSOStateInvalid(err error) bool {
	return errors.Is(err, errs.E(errs.CodeIdentitySSOStateInvalid))
}

var _ identityprovider.CallbackProvider = (*identityprovider.GitHubClient)(nil)
