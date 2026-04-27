package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewResetUserPasswordLogic creates a ResetUserPasswordLogic instance.
// NewResetUserPasswordLogic 创建 ResetUserPasswordLogic 实例。
func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	logCtx := withLogContext(ctx)
	return &ResetUserPasswordLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// ResetUserPassword adapts the gRPC request to the admin password service.
// ResetUserPassword 将 gRPC 请求适配到管理员重置密码 service。
func (l *ResetUserPasswordLogic) ResetUserPassword(in *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Services.Users.ResetUserPassword(l.ctx, identityservice.ResetUserPasswordInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		NewPassword:  in.GetNewPassword(),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "reset user password failed")
	}

	return &pb.ResetUserPasswordResponse{Ok: true}, nil
}
