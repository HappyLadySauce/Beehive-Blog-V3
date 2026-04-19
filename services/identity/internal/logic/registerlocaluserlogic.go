package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLocalUserLogic {
	return &RegisterLocalUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLocalUserLogic) RegisterLocalUser(in *pb.RegisterLocalUserRequest) (*pb.RegisterLocalUserResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RegisterLocalUserResponse{}, nil
}
