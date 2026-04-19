package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshSessionTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshSessionTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshSessionTokenLogic {
	return &RefreshSessionTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshSessionTokenLogic) RefreshSessionToken(in *pb.RefreshSessionTokenRequest) (*pb.RefreshSessionTokenResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RefreshSessionTokenResponse{}, nil
}
