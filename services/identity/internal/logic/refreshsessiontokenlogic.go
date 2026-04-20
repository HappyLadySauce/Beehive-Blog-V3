package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshSessionTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewRefreshSessionTokenLogic creates a RefreshSessionTokenLogic instance.
// NewRefreshSessionTokenLogic 创建 RefreshSessionTokenLogic 实例。
func NewRefreshSessionTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshSessionTokenLogic {
	return &RefreshSessionTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefreshSessionToken adapts the gRPC request to the refresh service.
// RefreshSessionToken 将 gRPC 请求适配到 refresh service。
func (l *RefreshSessionTokenLogic) RefreshSessionToken(in *pb.RefreshSessionTokenRequest) (*pb.RefreshSessionTokenResponse, error) {
	result, err := l.svcCtx.Services.Refresh.Execute(l.ctx, identityservice.RefreshSessionTokenInput{
		RefreshToken: in.GetRefreshToken(),
		UserAgent:    in.GetUserAgent(),
		ClientIP:     ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "refresh session token failed")
	}

	l.Infof("session refresh succeeded: user_id=%d session_id=%d", result.User.ID, result.Session.ID)

	return &pb.RefreshSessionTokenResponse{
		TokenPair: auth.NewTokenPair(
			result.AccessToken,
			result.RefreshToken,
			expiresIn(result.AccessExpiresAt),
			result.Session.ID,
		),
		SessionInfo: auth.ToSessionInfo(result.Session),
	}, nil
}
