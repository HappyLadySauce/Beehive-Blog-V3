package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewRegisterLocalUserLogic creates a RegisterLocalUserLogic instance.
// NewRegisterLocalUserLogic 创建 RegisterLocalUserLogic 实例。
func NewRegisterLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLocalUserLogic {
	return &RegisterLocalUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
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

	l.Infof("local registration succeeded: user_id=%d username=%s", result.User.ID, result.User.Username)

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
