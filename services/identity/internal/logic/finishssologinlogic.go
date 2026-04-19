package logic

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// FinishSsoLoginLogic handles SSO callback completion.
// FinishSsoLoginLogic 负责处理 SSO 回调完成逻辑。
type FinishSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewFinishSsoLoginLogic creates a FinishSsoLoginLogic instance.
// NewFinishSsoLoginLogic 创建 FinishSsoLoginLogic 实例。
func NewFinishSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinishSsoLoginLogic {
	return &FinishSsoLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FinishSsoLogin completes the SSO callback flow for fully implemented providers.
// FinishSsoLogin 完成已完整实现 provider 的 SSO 回调流程。
func (l *FinishSsoLoginLogic) FinishSsoLogin(in *pb.FinishSsoLoginRequest) (*pb.FinishSsoLoginResponse, error) {
	provider, err := auth.NormalizeProvider(in.GetProvider())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unsupported provider")
	}
	if provider != auth.ProviderGitHub {
		return nil, status.Error(codes.Unimplemented, "sso_provider_not_ready")
	}
	if strings.TrimSpace(in.GetCode()) == "" || strings.TrimSpace(in.GetState()) == "" {
		return nil, status.Error(codes.InvalidArgument, "code and state are required")
	}

	providerConf, err := auth.GetProviderConfig(l.svcCtx.Config.SSO, provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unsupported provider")
	}
	if !auth.IsProviderEnabled(l.svcCtx.Config.SSO, provider) {
		return nil, status.Error(codes.FailedPrecondition, "sso_provider_disabled")
	}
	if !auth.IsProviderLoginReady(provider) {
		return nil, status.Error(codes.FailedPrecondition, "sso_provider_not_ready")
	}
	redirectURI := strings.TrimSpace(in.GetRedirectUri())
	if redirectURI == "" || redirectURI != strings.TrimSpace(providerConf.RedirectURL) {
		return nil, status.Error(codes.InvalidArgument, "redirect_uri does not match configured provider redirect")
	}

	providerAccessToken, err := auth.ExchangeGitHubCode(l.ctx, providerConf, in.GetCode(), redirectURI)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "exchange github code failed: %v", err)
	}
	githubProfile, rawProfile, err := auth.FetchGitHubProfile(l.ctx, providerAccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "fetch github profile failed: %v", err)
	}

	now := time.Now().UTC()
	clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate refresh token failed: %v", err)
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	providerSubject := strconv.FormatInt(githubProfile.ID, 10)

	var user *entity.User
	var session *entity.UserSession

	err = l.svcCtx.Store.DB().WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		store := l.svcCtx.Store.WithTx(tx)

		stateRow, err := store.OAuthLoginStates.GetForUpdateByProviderState(l.ctx, provider, strings.TrimSpace(in.GetState()))
		if err != nil {
			if repo.IsNotFound(err) {
				return status.Error(codes.Unauthenticated, "sso_state_invalid")
			}
			return err
		}
		if stateRow.ConsumedAt != nil || stateRow.ExpiresAt.Before(now) {
			return status.Error(codes.Unauthenticated, "sso_state_invalid")
		}
		if stateRow.RedirectURI != redirectURI {
			return status.Error(codes.InvalidArgument, "redirect_uri mismatch")
		}

		fed, err := store.FederatedIdentities.GetByProviderSubject(l.ctx, provider, providerSubject)
		if err != nil && !repo.IsNotFound(err) {
			return err
		}

		if fed != nil {
			user, err = store.Users.GetByID(l.ctx, fed.UserID)
			if err != nil {
				return err
			}
			if err := validateActiveUserStatus(user.Status); err != nil {
				return err
			}
			if err := store.FederatedIdentities.TouchLogin(
				l.ctx,
				fed.ID,
				githubProfile.Email,
				stringPtr(strings.TrimSpace(githubProfile.Name)),
				githubProfile.AvatarURL,
				rawProfile,
				now,
			); err != nil {
				return err
			}
		} else {
			if githubProfile.Email != nil && *githubProfile.Email != "" {
				user, err = store.Users.GetByEmail(l.ctx, *githubProfile.Email)
				if err != nil && !repo.IsNotFound(err) {
					return err
				}
			}
			if user == nil {
				username, err := buildUniqueGitHubUsername(l.ctx, store, githubProfile.Login)
				if err != nil {
					return err
				}
				nickname := githubProfile.Name
				if strings.TrimSpace(nickname) == "" {
					nickname = githubProfile.Login
				}
				user = &entity.User{
					Username:    username,
					Email:       githubProfile.Email,
					Nickname:    stringPtr(strings.TrimSpace(nickname)),
					AvatarURL:   githubProfile.AvatarURL,
					Role:        auth.UserRoleMember,
					Status:      auth.UserStatusActive,
					LastLoginAt: &now,
				}
				if err := store.Users.Create(l.ctx, user); err != nil {
					return err
				}
			} else {
				if err := validateActiveUserStatus(user.Status); err != nil {
					return err
				}
			}

			providerLogin := githubProfile.Login
			fed = &entity.FederatedIdentity{
				UserID:              user.ID,
				Provider:            provider,
				ProviderSubject:     providerSubject,
				ProviderSubjectType: "github_user_id",
				ProviderLogin:       stringPtr(providerLogin),
				ProviderEmail:       githubProfile.Email,
				ProviderDisplayName: stringPtr(strings.TrimSpace(githubProfile.Name)),
				AvatarURL:           githubProfile.AvatarURL,
				AppIDOrClientID:     stringPtr(strings.TrimSpace(providerConf.ClientID)),
				AccessScope:         stringPtr(strings.Join(providerConf.Scopes, ",")),
				RawProfile:          rawProfile,
				LastLoginAt:         &now,
			}
			if err := store.FederatedIdentities.Create(l.ctx, fed); err != nil {
				return err
			}
		}

		if err := store.Users.TouchLogin(l.ctx, user.ID, now); err != nil {
			return err
		}

		session = &entity.UserSession{
			UserID:     user.ID,
			AuthSource: auth.AuthSourceSSO,
			ClientType: stringPtr(strings.TrimSpace(in.GetClientType())),
			DeviceID:   stringPtr(strings.TrimSpace(in.GetDeviceId())),
			DeviceName: stringPtr(strings.TrimSpace(in.GetDeviceName())),
			IPAddress:  stringPtr(clientIP),
			UserAgent:  stringPtr(strings.TrimSpace(in.GetUserAgent())),
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(l.svcCtx.Config.Security.RefreshTokenTTLSeconds) * time.Second),
		}
		if err := store.UserSessions.Create(l.ctx, session); err != nil {
			return err
		}

		refreshRow := &entity.RefreshToken{
			SessionID: session.ID,
			TokenHash: refreshTokenHash,
			IssuedAt:  now,
			ExpiresAt: session.ExpiresAt,
		}
		if err := store.RefreshTokens.Create(l.ctx, refreshRow); err != nil {
			return err
		}
		if err := store.OAuthLoginStates.Consume(l.ctx, stateRow.ID, now); err != nil {
			return err
		}

		authSource := auth.AuthSourceSSO
		writeAudit(l.ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &session.ID,
			Provider:   stringPtr(provider),
			AuthSource: &authSource,
			EventType:  auth.AuditEventFinishSSO,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(clientIP),
			UserAgent:  stringPtr(strings.TrimSpace(in.GetUserAgent())),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"provider_subject": providerSubject,
			}),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	accessToken, accessExpiresAt, err := issueAccessToken(
		l.svcCtx.Config.Security.AccessTokenSecret,
		l.svcCtx.Config.Security.AccessTokenTTLSeconds,
		user,
		session,
		now,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue access token failed: %v", err)
	}

	l.Infof("sso finish succeeded: provider=%s user_id=%d session_id=%d", provider, user.ID, session.ID)

	return &pb.FinishSsoLoginResponse{
		TokenPair: auth.NewTokenPair(
			accessToken,
			refreshToken,
			accessExpiresAt.Unix()-now.Unix(),
			session.ID,
		),
		CurrentUser: auth.ToCurrentUser(user),
		SessionInfo: auth.ToSessionInfo(session),
	}, nil
}
