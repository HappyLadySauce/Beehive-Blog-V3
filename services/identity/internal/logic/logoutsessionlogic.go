package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
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

// LogoutSession adapts the trusted logout request to the logout service.
// LogoutSession 将可信登出请求适配到 logout service。
func (l *LogoutSessionLogic) LogoutSession(in *pb.LogoutSessionRequest) (*pb.LogoutSessionResponse, error) {
	sessionID, err := parseID("session_id", in.GetSessionId())
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.Services.Logout.Execute(l.ctx, identityservice.LogoutSessionInput{
		SessionID: sessionID,
		ClientIP:  ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "logout session failed")
	}

	l.Infof("session logout succeeded: session_id=%d", sessionID)

	return &pb.LogoutSessionResponse{Ok: true}, nil
}
