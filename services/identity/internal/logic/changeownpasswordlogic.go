package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

type ChangeOwnPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewChangeOwnPasswordLogic creates a ChangeOwnPasswordLogic instance.
// NewChangeOwnPasswordLogic 创建 ChangeOwnPasswordLogic 实例。
func NewChangeOwnPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOwnPasswordLogic {
	logCtx := withLogContext(ctx)
	return &ChangeOwnPasswordLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// ChangeOwnPassword adapts the gRPC request to the password service.
// ChangeOwnPassword 将 gRPC 请求适配到密码修改 service。
func (l *ChangeOwnPasswordLogic) ChangeOwnPassword(in *pb.ChangeOwnPasswordRequest) (*pb.ChangeOwnPasswordResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.Services.Users.ChangeOwnPassword(l.ctx, identityservice.ChangeOwnPasswordInput{
		UserID:      userID,
		OldPassword: in.GetOldPassword(),
		NewPassword: in.GetNewPassword(),
		ClientIP:    ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	}); err != nil {
		return nil, toStatusError(err, "change own password failed")
	}

	return &pb.ChangeOwnPasswordResponse{Ok: true}, nil
}
