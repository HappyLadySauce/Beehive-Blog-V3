package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutSessionLogic {
	return &LogoutSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutSessionLogic) LogoutSession(in *pb.LogoutSessionRequest) (*pb.LogoutSessionResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.LogoutSessionResponse{}, nil
}
