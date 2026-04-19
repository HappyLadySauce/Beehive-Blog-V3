package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FinishSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFinishSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinishSsoLoginLogic {
	return &FinishSsoLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FinishSsoLoginLogic) FinishSsoLogin(in *pb.FinishSsoLoginRequest) (*pb.FinishSsoLoginResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.FinishSsoLoginResponse{}, nil
}
