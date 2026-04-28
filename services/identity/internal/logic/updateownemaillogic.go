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

// UpdateOwnEmailLogic handles current-user email updates.
// UpdateOwnEmailLogic 负责当前用户邮箱修改。
type UpdateOwnEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewUpdateOwnEmailLogic creates an UpdateOwnEmailLogic instance.
// NewUpdateOwnEmailLogic 创建 UpdateOwnEmailLogic 实例。
func NewUpdateOwnEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOwnEmailLogic {
	logCtx := withLogContext(ctx)
	return &UpdateOwnEmailLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// UpdateOwnEmail adapts current-user email update RPC requests to the service layer.
// UpdateOwnEmail 将当前用户邮箱修改 RPC 请求适配到 service 层。
func (l *UpdateOwnEmailLogic) UpdateOwnEmail(in *pb.UpdateOwnEmailRequest) (*pb.UpdateOwnEmailResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Services.Users.UpdateOwnEmail(l.ctx, identityservice.UpdateOwnEmailInput{
		UserID:             userID,
		Email:              in.GetEmail(),
		VerificationMethod: in.GetVerificationMethod(),
		CurrentPassword:    in.GetCurrentPassword(),
		Provider:           in.GetProvider(),
		Code:               in.GetCode(),
		State:              in.GetState(),
		RedirectURI:        in.GetRedirectUri(),
		ClientIP:           ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "update own email failed")
	}

	l.logger.Info("identity_update_own_email_succeeded", logs.Int64("user_id", userID))

	return &pb.UpdateOwnEmailResponse{CurrentUser: auth.ToCurrentUser(result.User)}, nil
}
