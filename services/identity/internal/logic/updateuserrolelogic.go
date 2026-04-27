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

type UpdateUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewUpdateUserRoleLogic creates an UpdateUserRoleLogic instance.
// NewUpdateUserRoleLogic 创建 UpdateUserRoleLogic 实例。
func NewUpdateUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRoleLogic {
	logCtx := withLogContext(ctx)
	return &UpdateUserRoleLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// UpdateUserRole adapts the gRPC request to the admin user service.
// UpdateUserRole 将 gRPC 请求适配到管理员用户 service。
func (l *UpdateUserRoleLogic) UpdateUserRole(in *pb.UpdateUserRoleRequest) (*pb.UpdateUserRoleResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateUserRole(l.ctx, identityservice.UpdateUserRoleInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Role:         roleString(in.GetRole()),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update user role failed")
	}

	return &pb.UpdateUserRoleResponse{User: auth.ToAdminUserView(result.User)}, nil
}
