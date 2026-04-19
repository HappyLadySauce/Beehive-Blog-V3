package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type IntrospectAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIntrospectAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntrospectAccessTokenLogic {
	return &IntrospectAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IntrospectAccessTokenLogic) IntrospectAccessToken(in *pb.IntrospectAccessTokenRequest) (*pb.IntrospectAccessTokenResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.IntrospectAccessTokenResponse{}, nil
}
