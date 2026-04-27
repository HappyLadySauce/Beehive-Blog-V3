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

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewUpdateUserStatusLogic creates an UpdateUserStatusLogic instance.
// NewUpdateUserStatusLogic 创建 UpdateUserStatusLogic 实例。
func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	logCtx := withLogContext(ctx)
	return &UpdateUserStatusLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// UpdateUserStatus adapts the gRPC request to the admin user service.
// UpdateUserStatus 将 gRPC 请求适配到管理员用户 service。
func (l *UpdateUserStatusLogic) UpdateUserStatus(in *pb.UpdateUserStatusRequest) (*pb.UpdateUserStatusResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateUserStatus(l.ctx, identityservice.UpdateUserStatusInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Status:       accountStatusString(in.GetStatus()),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update user status failed")
	}

	return &pb.UpdateUserStatusResponse{User: auth.ToAdminUserView(result.User)}, nil
}
