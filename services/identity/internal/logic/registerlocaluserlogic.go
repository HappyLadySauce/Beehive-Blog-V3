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

type RegisterLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewRegisterLocalUserLogic creates a RegisterLocalUserLogic instance.
// NewRegisterLocalUserLogic 创建 RegisterLocalUserLogic 实例。
func NewRegisterLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLocalUserLogic {
	logCtx := withLogContext(ctx)
	return &RegisterLocalUserLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// RegisterLocalUser adapts the gRPC request to the register service.
// RegisterLocalUser 将 gRPC 请求适配到注册 service。
func (l *RegisterLocalUserLogic) RegisterLocalUser(in *pb.RegisterLocalUserRequest) (*pb.RegisterLocalUserResponse, error) {
	result, err := l.svcCtx.Services.Register.Execute(l.ctx, identityservice.RegisterLocalUserInput{
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
		Nickname: in.GetNickname(),
		ClientIP: ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "register local user failed")
	}

	l.logger.Info(
		"identity_register_local_user_succeeded",
		logs.Int64("user_id", result.User.ID),
		logs.String("username", result.User.Username),
	)

	return &pb.RegisterLocalUserResponse{
		CurrentUser: auth.ToCurrentUser(result.User),
		TokenPair: auth.NewTokenPair(
			result.AccessToken,
			result.RefreshToken,
			expiresIn(result.AccessExpiresAt),
			result.Session.ID,
		),
		SessionInfo: auth.ToSessionInfo(result.Session),
	}, nil
}
