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

type UpdateUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewUpdateUserProfileLogic creates an UpdateUserProfileLogic instance.
// NewUpdateUserProfileLogic 创建 UpdateUserProfileLogic 实例。
func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	logCtx := withLogContext(ctx)
	return &UpdateUserProfileLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// UpdateUserProfile adapts the gRPC request to the admin user service.
// UpdateUserProfile 将 gRPC 请求适配到管理员用户资料 service。
func (l *UpdateUserProfileLogic) UpdateUserProfile(in *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	actorUserID, targetUserID, err := parseActorTarget(in.GetActorUserId(), in.GetTargetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateUserProfile(l.ctx, identityservice.UpdateUserProfileInput{
		ActorUserID:  actorUserID,
		TargetUserID: targetUserID,
		Username:     in.Username,
		Email:        in.Email,
		Nickname:     in.Nickname,
		AvatarURL:    in.AvatarUrl,
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update user profile failed")
	}

	return &pb.UpdateUserProfileResponse{User: auth.ToAdminUserView(result.User)}, nil
}
