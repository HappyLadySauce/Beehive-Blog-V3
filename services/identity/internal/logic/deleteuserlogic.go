package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	logCtx := withLogContext(ctx)
	return &DeleteUserLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// DeleteUser adapts the gRPC request to the admin soft-delete service.
// DeleteUser 将 gRPC 请求适配到管理员软删除 service。
func (l *DeleteUserLogic) DeleteUser(in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Services.Users.DeleteUser(l.ctx, identityservice.DeleteUserInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "delete user failed")
	}

	return &pb.DeleteUserResponse{Ok: true}, nil
}
