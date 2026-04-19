package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type LogoutSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewLogoutSessionLogic creates a LogoutSessionLogic instance.
// NewLogoutSessionLogic 创建 LogoutSessionLogic 实例。
func NewLogoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutSessionLogic {
	return &LogoutSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LogoutSession revokes the current trusted session and its refresh tokens.
// LogoutSession 吊销当前可信会话及其 refresh token。
func (l *LogoutSessionLogic) LogoutSession(in *pb.LogoutSessionRequest) (*pb.LogoutSessionResponse, error) {
	// Parse and validate the trusted session identifier.
	// 解析并校验可信会话标识。
	sessionID, err := parseID("session_id", in.GetSessionId())
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)

	err = l.svcCtx.Store.DB().WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		store := l.svcCtx.Store.WithTx(tx)
		session, err := store.UserSessions.GetForUpdateByID(l.ctx, sessionID)
		if err != nil {
			if repo.IsNotFound(err) {
				return status.Error(codes.NotFound, "session not found")
			}
			return err
		}

		if session.Status == auth.SessionStatusActive {
			if err := store.UserSessions.Revoke(l.ctx, session.ID, now); err != nil {
				return err
			}
		}
		if err := store.RefreshTokens.RevokeActiveBySessionID(l.ctx, session.ID, now); err != nil {
			return err
		}

		authSource := session.AuthSource
		writeAudit(l.ctx, store, auditInput{
			UserID:     &session.UserID,
			SessionID:  &session.ID,
			AuthSource: &authSource,
			EventType:  auth.AuditEventLogoutSession,
			Result:     auth.AuditResultSuccess,
			ClientIP:   stringPtr(clientIP),
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	l.Infof("session logout succeeded: session_id=%d", sessionID)

	return &pb.LogoutSessionResponse{Ok: true}, nil
}
