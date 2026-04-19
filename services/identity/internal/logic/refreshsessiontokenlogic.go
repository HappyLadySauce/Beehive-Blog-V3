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

type RefreshSessionTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewRefreshSessionTokenLogic creates a RefreshSessionTokenLogic instance.
// NewRefreshSessionTokenLogic 创建 RefreshSessionTokenLogic 实例。
func NewRefreshSessionTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshSessionTokenLogic {
	return &RefreshSessionTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefreshSessionToken rotates the refresh token and reissues the access token.
// RefreshSessionToken 轮换 refresh token 并重新签发 access token。
func (l *RefreshSessionTokenLogic) RefreshSessionToken(in *pb.RefreshSessionTokenRequest) (*pb.RefreshSessionTokenResponse, error) {
	// Validate refresh token input before any persistence operation.
	// 在执行任何持久化操作前校验 refresh token 输入。
	rawRefreshToken := strings.TrimSpace(in.GetRefreshToken())
	if rawRefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	now := time.Now().UTC()
	refreshTokenHash := auth.HashRefreshToken(rawRefreshToken)
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate refresh token failed: %v", err)
	}
	newRefreshTokenHash := auth.HashRefreshToken(newRefreshToken)
	clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)

	var user *entity.User
	var session *entity.UserSession

	err = l.svcCtx.Store.DB().WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		store := l.svcCtx.Store.WithTx(tx)

		currentToken, err := store.RefreshTokens.GetActiveForUpdateByHash(l.ctx, refreshTokenHash)
		if err != nil {
			if repo.IsNotFound(err) {
				return status.Error(codes.Unauthenticated, "invalid_refresh_token")
			}
			return err
		}
		if currentToken.ExpiresAt.Before(now) {
			_ = store.RefreshTokens.Revoke(l.ctx, currentToken.ID, now)
			return status.Error(codes.Unauthenticated, "refresh_token_expired")
		}

		session, err = store.UserSessions.GetForUpdateByID(l.ctx, currentToken.SessionID)
		if err != nil {
			if repo.IsNotFound(err) {
				return status.Error(codes.Unauthenticated, "session_revoked")
			}
			return err
		}
		if session.Status != auth.SessionStatusActive {
			return status.Error(codes.Unauthenticated, "session_revoked")
		}
		if session.ExpiresAt.Before(now) {
			_ = store.UserSessions.MarkExpired(l.ctx, session.ID, now)
			return status.Error(codes.Unauthenticated, "session_expired")
		}

		user, err = store.Users.GetByID(l.ctx, session.UserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return status.Error(codes.Unauthenticated, "account_not_found")
			}
			return err
		}
		if err := validateActiveUserStatus(user.Status); err != nil {
			return err
		}

		session.LastSeenAt = &now
		session.ExpiresAt = now.Add(time.Duration(l.svcCtx.Config.Security.RefreshTokenTTLSeconds) * time.Second)
		if err := store.UserSessions.TouchActive(l.ctx, session.ID, now, session.ExpiresAt); err != nil {
			return err
		}
		if err := store.RefreshTokens.Revoke(l.ctx, currentToken.ID, now); err != nil {
			return err
		}

		newTokenRow := &entity.RefreshToken{
			SessionID:          session.ID,
			TokenHash:          newRefreshTokenHash,
			IssuedAt:           now,
			ExpiresAt:          session.ExpiresAt,
			RotatedFromTokenID: &currentToken.ID,
		}
		if err := store.RefreshTokens.Create(l.ctx, newTokenRow); err != nil {
			return err
		}

		authSource := session.AuthSource
		writeAudit(l.ctx, store, auditInput{
			UserID:     &user.ID,
			SessionID:  &session.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventRefreshSession,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(clientIP),
			UserAgent:  stringPtr(strings.TrimSpace(in.GetUserAgent())),
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

	l.Infof("session refresh succeeded: user_id=%d session_id=%d", user.ID, session.ID)

	return &pb.RefreshSessionTokenResponse{
		TokenPair: auth.NewTokenPair(
			accessToken,
			newRefreshToken,
			accessExpiresAt.Unix()-now.Unix(),
			session.ID,
		),
		SessionInfo: auth.ToSessionInfo(session),
	}, nil
}
