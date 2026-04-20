package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type LoginLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewLoginLocalUserLogic creates a LoginLocalUserLogic instance.
// NewLoginLocalUserLogic 创建 LoginLocalUserLogic 实例。
func NewLoginLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLocalUserLogic {
	logCtx := withLogContext(ctx)
	return &LoginLocalUserLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// LoginLocalUser adapts the gRPC request to the login service.
// LoginLocalUser 将 gRPC 请求适配到登录 service。
func (l *LoginLocalUserLogic) LoginLocalUser(in *pb.LoginLocalUserRequest) (*pb.LoginLocalUserResponse, error) {
	result, err := l.svcCtx.Services.Login.Execute(l.ctx, identityservice.LoginLocalUserInput{
		LoginIdentifier: in.GetLoginIdentifier(),
		Password:        in.GetPassword(),
		ClientType:      in.GetClientType(),
		DeviceID:        in.GetDeviceId(),
		DeviceName:      in.GetDeviceName(),
		UserAgent:       in.GetUserAgent(),
		ClientIP:        ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "login local user failed")
	}

	l.logger.Info(
		"identity_login_local_user_succeeded",
		logs.Int64("user_id", result.User.ID),
		logs.Int64("session_id", result.Session.ID),
	)

	return &pb.LoginLocalUserResponse{
		TokenPair: auth.NewTokenPair(
			result.AccessToken,
			result.RefreshToken,
			expiresIn(result.AccessExpiresAt),
			result.Session.ID,
		),
		CurrentUser: auth.ToCurrentUser(result.User),
		SessionInfo: auth.ToSessionInfo(result.Session),
	}, nil
}
