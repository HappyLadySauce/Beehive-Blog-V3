package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type GetCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewGetCurrentUserLogic creates a GetCurrentUserLogic instance.
// NewGetCurrentUserLogic 创建 GetCurrentUserLogic 实例。
func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	logCtx := withLogContext(ctx)
	return &GetCurrentUserLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// GetCurrentUser adapts the trusted identity request to the current-user service.
// GetCurrentUser 将可信身份请求适配到 current-user service。
func (l *GetCurrentUserLogic) GetCurrentUser(in *pb.GetCurrentUserRequest) (*pb.GetCurrentUserResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Services.CurrentUser.Execute(l.ctx, identityservice.GetCurrentUserInput{
		UserID: userID,
	})
	if err != nil {
		return nil, toStatusError(err, "get current user failed")
	}

	return &pb.GetCurrentUserResponse{
		CurrentUser: auth.ToCurrentUser(result.User),
	}, nil
}
