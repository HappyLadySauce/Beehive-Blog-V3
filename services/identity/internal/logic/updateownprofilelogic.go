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

type UpdateOwnProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewUpdateOwnProfileLogic creates an UpdateOwnProfileLogic instance.
// NewUpdateOwnProfileLogic 创建 UpdateOwnProfileLogic 实例。
func NewUpdateOwnProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOwnProfileLogic {
	logCtx := withLogContext(ctx)
	return &UpdateOwnProfileLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// UpdateOwnProfile adapts the gRPC request to the profile service.
// UpdateOwnProfile 将 gRPC 请求适配到资料更新 service。
func (l *UpdateOwnProfileLogic) UpdateOwnProfile(in *pb.UpdateOwnProfileRequest) (*pb.UpdateOwnProfileResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}
	result, err := l.svcCtx.Services.Users.UpdateOwnProfile(l.ctx, identityservice.UpdateOwnProfileInput{
		UserID:    userID,
		Nickname:  in.Nickname,
		AvatarURL: in.AvatarUrl,
		ClientIP:  ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update own profile failed")
	}

	return &pb.UpdateOwnProfileResponse{CurrentUser: auth.ToCurrentUser(result.User)}, nil
}
