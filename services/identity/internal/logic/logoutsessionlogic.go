package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type LogoutSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewLogoutSessionLogic creates a LogoutSessionLogic instance.
// NewLogoutSessionLogic 创建 LogoutSessionLogic 实例。
func NewLogoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutSessionLogic {
	logCtx := withLogContext(ctx)
	return &LogoutSessionLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
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

	l.logger.Info(
		"identity_logout_session_succeeded",
		logs.Int64("session_id", sessionID),
	)

	return &pb.LogoutSessionResponse{Ok: true}, nil
}
