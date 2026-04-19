package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSsoLoginLogic {
	return &StartSsoLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StartSsoLoginLogic) StartSsoLogin(in *pb.StartSsoLoginRequest) (*pb.StartSsoLoginResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.StartSsoLoginResponse{}, nil
}
