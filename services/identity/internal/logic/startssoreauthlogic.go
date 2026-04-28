package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// StartSsoReauthLogic handles SSO reauthentication URL generation.
// StartSsoReauthLogic 负责生成 SSO 重验授权地址。
type StartSsoReauthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewStartSsoReauthLogic creates a StartSsoReauthLogic instance.
// NewStartSsoReauthLogic 创建 StartSsoReauthLogic 实例。
func NewStartSsoReauthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSsoReauthLogic {
	logCtx := withLogContext(ctx)
	return &StartSsoReauthLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// StartSsoReauth adapts SSO reauthentication start requests to the service layer.
// StartSsoReauth 将 SSO 重验发起请求适配到 service 层。
func (l *StartSsoReauthLogic) StartSsoReauth(in *pb.StartSsoReauthRequest) (*pb.StartSsoLoginResponse, error) {
	userID, err := parseID("user_id", in.GetUserId())
	if err != nil {
		return nil, err
	}

	result, err := l.svcCtx.Services.SSOStart.ExecuteReauth(l.ctx, identityservice.StartSSOReauthInput{
		UserID:      userID,
		Provider:    in.GetProvider(),
		RedirectURI: in.GetRedirectUri(),
		State:       in.GetState(),
		ClientIP:    ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "start sso reauth failed")
	}

	l.logger.Info("identity_start_sso_reauth_succeeded", logs.String("provider", result.Provider))

	return &pb.StartSsoLoginResponse{
		Provider: result.Provider,
		AuthUrl:  result.AuthURL,
		State:    result.State,
	}, nil
}
