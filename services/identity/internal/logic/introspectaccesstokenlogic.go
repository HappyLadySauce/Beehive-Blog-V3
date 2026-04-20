package logic

import (
	"context"
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IntrospectAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewIntrospectAccessTokenLogic creates an IntrospectAccessTokenLogic instance.
// NewIntrospectAccessTokenLogic 创建 IntrospectAccessTokenLogic 实例。
func NewIntrospectAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntrospectAccessTokenLogic {
	return &IntrospectAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// IntrospectAccessToken adapts token introspection to the service layer.
// IntrospectAccessToken 将 token introspection 适配到 service 层。
func (l *IntrospectAccessTokenLogic) IntrospectAccessToken(in *pb.IntrospectAccessTokenRequest) (*pb.IntrospectAccessTokenResponse, error) {
	result, err := l.svcCtx.Services.Introspect.Execute(l.ctx, identityservice.IntrospectAccessTokenInput{
		AccessToken: in.GetAccessToken(),
	})
	if err != nil {
		return nil, toStatusError(err, "introspect access token failed")
	}
	if !result.Active {
		return &pb.IntrospectAccessTokenResponse{Active: false}, nil
	}

	return &pb.IntrospectAccessTokenResponse{
		Active:        true,
		UserId:        strconv.FormatInt(result.User.ID, 10),
		Role:          auth.ToProtoRole(result.User.Role),
		AccountStatus: auth.ToProtoAccountStatus(result.User.Status),
		SessionId:     strconv.FormatInt(result.Session.ID, 10),
		AuthSource:    auth.ToProtoAuthSource(result.Session.AuthSource),
		ExpiresAt:     result.ExpiresAt,
	}, nil
}
