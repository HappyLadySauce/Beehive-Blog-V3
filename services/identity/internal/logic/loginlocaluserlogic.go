package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLocalUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLocalUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLocalUserLogic {
	return &LoginLocalUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLocalUserLogic) LoginLocalUser(in *pb.LoginLocalUserRequest) (*pb.LoginLocalUserResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginLocalUserResponse{}, nil
}
