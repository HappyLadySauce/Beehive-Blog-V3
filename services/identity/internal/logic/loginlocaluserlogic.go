package logic

import (
	"context"
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

type LoginLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewLoginLocalUserLogic creates a LoginLocalUserLogic instance.
// NewLoginLocalUserLogic 创建 LoginLocalUserLogic 实例。
func NewLoginLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLocalUserLogic {
	return &LoginLocalUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LoginLocalUser authenticates a local user and creates a new session.
// LoginLocalUser 认证本地用户并创建新会话。
func (l *LoginLocalUserLogic) LoginLocalUser(in *pb.LoginLocalUserRequest) (*pb.LoginLocalUserResponse, error) {
	// Normalize and validate the login request.
	// 规范化并校验登录请求。
	identifier, err := auth.NormalizeLoginIdentifier(in.GetLoginIdentifier())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if strings.TrimSpace(in.GetPassword()) == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	user, err := l.svcCtx.Store.Users.FindByLoginIdentifier(l.ctx, identifier)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid_credentials")
		}
		return nil, status.Errorf(codes.Internal, "lookup user failed: %v", err)
	}
	if err := validateActiveUserStatus(user.Status); err != nil {
		return nil, err
	}

	credential, err := l.svcCtx.Store.CredentialLocals.GetByUserID(l.ctx, user.ID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid_credentials")
		}
		return nil, status.Errorf(codes.Internal, "lookup credential failed: %v", err)
	}
	if err := auth.VerifyPassword(credential.PasswordHash, in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid_credentials")
	}

	now := time.Now().UTC()
	clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)
	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate refresh token failed: %v", err)
	}
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	authSource := auth.AuthSourceLocal

	var createdSession *entity.UserSession

	err = l.svcCtx.Store.DB().WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		store := l.svcCtx.Store.WithTx(tx)
		if err := store.Users.TouchLogin(l.ctx, user.ID, now); err != nil {
			return err
		}

		createdSession = &entity.UserSession{
			UserID:     user.ID,
			AuthSource: auth.AuthSourceLocal,
			ClientType: stringPtr(strings.TrimSpace(in.GetClientType())),
			DeviceID:   stringPtr(strings.TrimSpace(in.GetDeviceId())),
			DeviceName: stringPtr(strings.TrimSpace(in.GetDeviceName())),
			IPAddress:  stringPtr(clientIP),
			UserAgent:  stringPtr(strings.TrimSpace(in.GetUserAgent())),
			Status:     auth.SessionStatusActive,
			LastSeenAt: &now,
			ExpiresAt:  now.Add(time.Duration(l.svcCtx.Config.Security.RefreshTokenTTLSeconds) * time.Second),
		}
		if err := store.UserSessions.Create(l.ctx, createdSession); err != nil {
			return err
		}

		refreshRow := &entity.RefreshToken{
			SessionID: createdSession.ID,
			TokenHash: refreshTokenHash,
			IssuedAt:  now,
			ExpiresAt: createdSession.ExpiresAt,
		}
		if err := store.RefreshTokens.Create(l.ctx, refreshRow); err != nil {
			return err
		}

		writeAudit(l.ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &createdSession.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventLoginLocal,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(clientIP),
			UserAgent:  stringPtr(strings.TrimSpace(in.GetUserAgent())),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"login_identifier": identifier,
			}),
		})

		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login local user failed: %v", err)
	}

	user.LastLoginAt = &now
	accessToken, accessExpiresAt, err := issueAccessToken(
		l.svcCtx.Config.Security.AccessTokenSecret,
		l.svcCtx.Config.Security.AccessTokenTTLSeconds,
		user,
		createdSession,
		now,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "issue access token failed: %v", err)
	}

	l.Infof("local login succeeded: user_id=%d session_id=%d", user.ID, createdSession.ID)

	return &pb.LoginLocalUserResponse{
		TokenPair: auth.NewTokenPair(
			accessToken,
			refreshToken,
			accessExpiresAt.Unix()-now.Unix(),
			createdSession.ID,
		),
		CurrentUser: auth.ToCurrentUser(user),
		SessionInfo: auth.ToSessionInfo(createdSession),
	}, nil
}
